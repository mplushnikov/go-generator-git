package implementation

import (
	"context"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
)

type GitGeneratorImpl struct {
}

func (g *GitGeneratorImpl) CreateTemporaryWorkdir(ctx context.Context, basePath string) error {
	panic("implement me")
}

func (g *GitGeneratorImpl) CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string) error {
	panic("implement me")
}

func (g *GitGeneratorImpl) CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string) error {
	panic("implement me")
}

func (g *GitGeneratorImpl) WriteRenderSpecFile(ctx context.Context, generatorName string, renderSpecFile string, parameters map[string]string) (*genlibapi.Response, error) {
	panic("implement me")
}

func (g *GitGeneratorImpl) Generate(ctx context.Context) (*genlibapi.Response, error) {
	panic("implement me")
}

func (g *GitGeneratorImpl) Cleanup(ctx context.Context) error {
	panic("implement me")
}
