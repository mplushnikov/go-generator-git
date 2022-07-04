package generatorgit

import (
	"context"
	"github.com/StephanHCB/go-generator-git/v2/api"
	"github.com/StephanHCB/go-generator-git/v2/internal/implementation"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

var Instance api.GitApi

func init() {
	Instance = &implementation.GitGeneratorImpl{}
}

func ThreadsafeInstance() api.GitApi {
	return &implementation.GitGeneratorImpl{}
}

// this is not thread safe - only use for tests or cmd line client

func CreateTemporaryWorkdir(ctx context.Context, basePath string) error {
	return Instance.CreateTemporaryWorkdir(ctx, basePath)
}

func CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	return Instance.CloneSourceRepo(ctx, gitRepoUrl, gitBranch, auth)
}

func CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	return Instance.CloneTargetRepo(ctx, gitRepoUrl, gitBranch, baseBranch, auth)
}

func PrepareTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (api.GitApiRepo, error) {
	return Instance.PrepareTargetRepo(ctx, gitRepoUrl, gitBranch, auth)
}

func WriteRenderSpecFile(ctx context.Context, generatorName string, renderSpecFile string, parameters map[string]interface{}) (*genlibapi.Response, error) {
	return Instance.WriteRenderSpecFile(ctx, generatorName, renderSpecFile, parameters)
}

func Generate(ctx context.Context) (*genlibapi.Response, error) {
	return Instance.Generate(ctx)
}

func CommitAndPush(ctx context.Context, name string, email string, message string, auth transport.AuthMethod) error {
	return Instance.CommitAndPush(ctx, name, email, message, auth)
}

func Cleanup(ctx context.Context) error {
	return Instance.Cleanup(ctx)
}
