package api

import (
	"context"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitApiRepo interface {
	GetLocalPath() string
}

// Functionality that this library exposes.
type GitApi interface {
	// create a temporary working directory with a random name underneath basePath
	//
	// You need to call this before any of the other methods can be called
	//
	// We use a random sub directory so multiple goroutines can render in parallel
	CreateTemporaryWorkdir(ctx context.Context, basePath string) error

	// clone the source repo into the working directory and switch to the given branch (or tag, or revision)
	CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// clone the target repo into the working directory and set up the given branch
	//
	// if the branch does not yet exist, it will be created from the base branch (or tag, or revision),
	// otherwise we just check it out.
	CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// prepare the target repo into the working directory
	PrepareTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// write the given parameters for the given generator to a render spec file in the target directory
	//
	// unless some specific reason prevents you from this naming convention, renderSpecFile should be
	// 'generated-<generatorName>.yaml'
	//
	// the parameters will be validated against the generator spec file found in the source repo (called
	// 'generator-<generatorName>.yaml'). It is an error if any parameter does not conform to the specification,
	// or is missing and does not have a default, or if any parameters are unknown.
	//
	// Note that the render spec file in the target directory is silently overwritten. It is a git repo after all.
	// If the render spec file does not exist, that just means you are using the generator for the first time,
	// so it is silently created.
	//
	// Response is filled even in case of an error and will contain more details of what caused the error.
	WriteRenderSpecFile(ctx context.Context,
		generatorName string,
		renderSpecFile string,
		parameters map[string]interface{}) (*genlibapi.Response, error)

	// generate files using the render spec file written by WriteRenderSpecFile
	//
	// Response is filled even in case of an error and will contain more details of what caused the error
	// and what output files were affected. After a successful run, Response also contains the list of files
	// that were rendered.
	Generate(ctx context.Context) (*genlibapi.Response, error)

	// commit the changes in the target and push them (if an auth method is supplied)
	CommitAndPush(ctx context.Context, name string, email string, message string, auth transport.AuthMethod) error

	// delete the temporary working directory, including the source and target clones underneath it
	//
	// the base path given to CreateTemporaryWorkdir is left untouched so it can be re-used for the next
	// (or concurrent) render operations.
	Cleanup(ctx context.Context) error
}
