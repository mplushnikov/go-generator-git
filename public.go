package generatorgit

import (
	"context"
	"github.com/StephanHCB/go-generator-git/api"
	"github.com/StephanHCB/go-generator-git/internal/implementation"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
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

func CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string) error {
	return Instance.CloneSourceRepo(ctx, gitRepoUrl, gitBranch)
}

func CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string) error {
	return Instance.CloneTargetRepo(ctx, gitRepoUrl, gitBranch, baseBranch)
}

func WriteRenderSpecFile(ctx context.Context, generatorName string, renderSpecFile string, parameters map[string]string) (*genlibapi.Response, error) {
	return Instance.WriteRenderSpecFile(ctx, generatorName, renderSpecFile, parameters)
}

func Generate(ctx context.Context) (*genlibapi.Response, error) {
	return Instance.Generate(ctx)
}

func CommitAndPush(ctx context.Context, name string, email string, message string) error {
	return Instance.CommitAndPush(ctx, name, email, message)
}

func Cleanup(ctx context.Context) error {
	return Instance.Cleanup(ctx)
}
