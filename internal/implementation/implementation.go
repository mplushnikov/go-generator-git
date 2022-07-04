package implementation

import (
	"context"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-generator-git/v2/api"
	"github.com/StephanHCB/go-generator-git/v2/internal/repository/gitsourcerepo"
	"github.com/StephanHCB/go-generator-git/v2/internal/repository/gittargetrepo"
	"github.com/StephanHCB/go-generator-git/v2/internal/repository/tmpdir"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"path/filepath"
)

type GitGeneratorImpl struct {
	workdir        *tmpdir.TmpDir
	source         *gitsourcerepo.GitSourceRepo
	target         *gittargetrepo.GitTargetRepo
	targetBranch   string
	renderSpecFile string
}

type GitApiRepoImpl struct {
	localPath string
}

func (this *GitApiRepoImpl) GetLocalPath() string {
	return this.localPath
}

func (g *GitGeneratorImpl) CreateTemporaryWorkdir(ctx context.Context, basePath string) error {
	g.workdir = tmpdir.Instance(ctx, basePath)
	aulogging.Logger.Ctx(ctx).Debug().Printf("creating temporary working directory %s", g.workdir.Path(ctx))
	return g.workdir.Create(ctx)
}

func (g *GitGeneratorImpl) CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	if g.workdir == nil {
		return nil, errCreateWorkdirFirst(ctx)
	}
	if g.source != nil {
		return nil, errDuplicateClone(ctx, "source")
	}
	path := filepath.Join(g.workdir.Path(ctx), "source")
	aulogging.Logger.Ctx(ctx).Info().Printf("cloning source repo to %s", path)
	g.source = gitsourcerepo.Instance(ctx, path)
	if err := g.source.Clone(ctx, gitRepoUrl, gitBranch, auth); err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error cloning source repo from %s on branch %s", gitRepoUrl, gitBranch)
		return &GitApiRepoImpl{path}, err
	}
	return &GitApiRepoImpl{path}, nil
}

func (g *GitGeneratorImpl) PrepareTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	if g.workdir == nil {
		return nil, errCreateWorkdirFirst(ctx)
	}
	if g.target != nil {
		return nil, errDuplicateClone(ctx, "target")
	}

	path := filepath.Join(g.workdir.Path(ctx), "target")

	g.target = gittargetrepo.Instance(ctx, path)
	err := g.target.PrepareInit(ctx, gitRepoUrl, gitBranch)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error preparing target repo from %s", gitRepoUrl)
		return &GitApiRepoImpl{path}, err
	}
	g.targetBranch = gitBranch

	return &GitApiRepoImpl{path}, nil
}

func (g *GitGeneratorImpl) CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	if g.workdir == nil {
		return nil, errCreateWorkdirFirst(ctx)
	}
	if g.target != nil {
		return nil, errDuplicateClone(ctx, "target")
	}

	path := filepath.Join(g.workdir.Path(ctx), "target")
	localApiRepo := &GitApiRepoImpl{path}

	aulogging.Logger.Ctx(ctx).Info().Printf("cloning target repo to %s", path)
	g.target = gittargetrepo.Instance(ctx, path)
	if err := g.target.Clone(ctx, gitRepoUrl, auth); err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error cloning target repo from %s", gitRepoUrl)
		return localApiRepo, err
	}

	if hash := g.target.GetHashForRevision(ctx, gitBranch); hash != nil {
		aulogging.Logger.Ctx(ctx).Info().Printf("checking out %s (currently at %s)", gitBranch, hash.String())
		if err := g.target.Checkout(ctx, gitBranch); err != nil {
			aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error checking out %s", gitBranch)
			return localApiRepo, err
		} else {
			// ok. remember it for CommitAndPush()
			g.targetBranch = gitBranch
			return localApiRepo, nil
		}
	} else {
		if baseHash := g.target.GetHashForRevision(ctx, baseBranch); baseHash != nil {
			aulogging.Logger.Ctx(ctx).Debug().Printf("target branch %s not found - will create it", gitBranch)
			aulogging.Logger.Ctx(ctx).Debug().Printf("base branch %s is at %s - starting from there", baseBranch, baseHash.String())

			aulogging.Logger.Ctx(ctx).Info().Printf("creating new branch %s from %s", gitBranch, baseHash.String())
			if err := g.target.CreateBranch(ctx, gitBranch, baseHash); err != nil {
				aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error creating branch %s from %s", gitBranch, baseHash.String())
				return localApiRepo, err
			}

			// now the branch exists, verify and check it out
			hash = g.target.GetHashForRevision(ctx, gitBranch)
			if hash == nil {
				message := "internal error - lookup of branch failed right after create"
				aulogging.Logger.Ctx(ctx).Error().Print(message)
				return localApiRepo, errors.New(message)
			}

			aulogging.Logger.Ctx(ctx).Info().Printf("now checking out %s (currently at %s)", gitBranch, hash.String())
			if err := g.target.Checkout(ctx, gitBranch); err != nil {
				aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("error checking out %s", gitBranch)
				return localApiRepo, err
			}

			// ok. remember it for CommitAndPush()
			g.targetBranch = gitBranch
			return localApiRepo, nil
		} else {
			message := fmt.Sprintf("base branch %s does not exist", baseBranch)
			aulogging.Logger.Ctx(ctx).Error().Print(message)
			return localApiRepo, errors.New(message)
		}
	}
}

func (g *GitGeneratorImpl) WriteRenderSpecFile(ctx context.Context, generatorName string, renderSpecFile string, parameters map[string]interface{}) (*genlibapi.Response, error) {
	if g.workdir == nil {
		return &genlibapi.Response{Success: false}, errCreateWorkdirFirst(ctx)
	}
	if g.source == nil {
		return &genlibapi.Response{Success: false}, errCloneSourceFirst(ctx)
	}
	if g.target == nil {
		return &genlibapi.Response{Success: false}, errCloneTargetFirst(ctx)
	}
	if g.targetBranch == "" {
		return &genlibapi.Response{Success: false}, errCloneTargetSuccessfullyFirst(ctx)
	}

	// set it for request() and remember it for Generate()
	g.renderSpecFile = renderSpecFile

	response := generatorlib.WriteRenderSpecWithValues(ctx, g.request(), generatorName, parameters)
	if !response.Success {
		return response, errors.New("writing render spec file failed, see response for details")
	}
	return response, nil
}

func (g *GitGeneratorImpl) Generate(ctx context.Context) (*genlibapi.Response, error) {
	if g.workdir == nil {
		return &genlibapi.Response{Success: false}, errCreateWorkdirFirst(ctx)
	}
	if g.source == nil {
		return &genlibapi.Response{Success: false}, errCloneSourceFirst(ctx)
	}
	if g.target == nil {
		return &genlibapi.Response{Success: false}, errCloneTargetFirst(ctx)
	}
	if g.targetBranch == "" {
		return &genlibapi.Response{Success: false}, errCloneTargetSuccessfullyFirst(ctx)
	}
	if g.renderSpecFile == "" {
		return &genlibapi.Response{Success: false}, errWriteRenderSpecFirst(ctx)
	}

	response := generatorlib.Render(ctx, g.request())
	if !response.Success {
		return response, errors.New("rendering failed, see response for details")
	}
	return response, nil
}

func (g *GitGeneratorImpl) CommitAndPush(ctx context.Context, name string, email string, message string, auth transport.AuthMethod) error {
	if g.workdir == nil {
		return errCreateWorkdirFirst(ctx)
	}
	if g.target == nil {
		return errCloneTargetFirst(ctx)
	}
	if g.targetBranch == "" {
		return errCloneTargetSuccessfullyFirst(ctx)
	}

	if auth != nil {
		g.target.EnablePush()
		aulogging.Logger.Ctx(ctx).Info().Printf("committing and pushing")
	} else {
		aulogging.Logger.Ctx(ctx).Info().Printf("committing (but not pushing - no auth supplied)")
	}
	err := g.target.CommitAndPush(ctx, name, email, message, auth)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("error during commit or push")
		return err
	}

	return nil
}

func (g *GitGeneratorImpl) Cleanup(ctx context.Context) error {
	if g.workdir == nil {
		aulogging.Logger.Ctx(ctx).Debug().Print("skipping cleanup of temporary working directory that was never created")
		return nil
	}
	aulogging.Logger.Ctx(ctx).Debug().Printf("cleaning up temporary working directory %s", g.workdir.Path(ctx))
	return g.workdir.DeleteRecursive(ctx)
}

// internals

func (g *GitGeneratorImpl) request() *genlibapi.Request {
	return &genlibapi.Request{
		SourceBaseDir:  g.source.Path(),
		TargetBaseDir:  g.target.Path(),
		RenderSpecFile: g.renderSpecFile,
	}
}

// error situations

func errCreateWorkdirFirst(ctx context.Context) error {
	return errMsg(ctx, "implementation error - need to create a temporary workdir before clone")
}

func errDuplicateClone(ctx context.Context, whichRepo string) error {
	return errMsg(ctx, "implementation error - duplicate clone for "+whichRepo)
}

func errCloneSourceFirst(ctx context.Context) error {
	return errMsg(ctx, "implementation error - must clone source before using it")
}

func errCloneTargetFirst(ctx context.Context) error {
	return errMsg(ctx, "implementation error - must clone target before making changes to it")
}

func errCloneTargetSuccessfullyFirst(ctx context.Context) error {
	return errMsg(ctx, "implementation error - target clone or branch checkout was not successful, you cannot make changes to it")
}

func errWriteRenderSpecFirst(ctx context.Context) error {
	return errMsg(ctx, "implementation error - you must write the render spec file before templates can be rendered")
}

func errMsg(ctx context.Context, message string) error {
	aulogging.Logger.Ctx(ctx).Error().Print(message)
	return errors.New(message)
}
