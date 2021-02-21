# go-generator-git

A golang library for generating files from templates, capable of obtaining the templates from one git repository
and committing and pushing the result to another git repository. Can be used to scaffold application code.

See [go-generator-lib](https://github.com/StephanHCB/go-generator-lib/) for details regarding the templates and
the generation process. This library just clones the git repositories to two directories, 
creates or checks out the necessary branches, and handles the final commit and push in the target repository. 
It relies on _go-generator-lib_ to generate output files between the two directories.

## Usage

### Globals (not thread safe)

For easy use in one-off situations or command line tools, we have added global accessor functions.

**important: these operate on a global singleton instance, and thus are not thread safe** 

Leaving out any error handling, this is the minimal code to perform a full clone, render, commit cycle:

```
parameters := map[string]string{} // assuming all parameters have defaults, otherwise specify here

aulogging.SetupNoLoggerForTesting()

generatorgit.CreateTemporaryWorkdir(context.TODO(), "../output")
defer generatorgit.Cleanup(context.TODO())

generatorgit.CloneSourceRepo(context.TODO(), "https://github.com/StephanHCB/tpl-go-rest-chi", "master")
generatorgit.CloneTargetRepo(context.TODO(), "https://github.com/StephanHCB/scratch", "feature/target", "main")
generatorgit.WriteRenderSpecFile(context.TODO(), "main", "generated-main.yaml", parameters)
generatorgit.Generate(context.TODO())
generatorgit.CommitAndPush(context.TODO(), "somebody", "somebody@mailinator.com", "initial generation", nil)
```

Note that `CommitAndPush` will only push the commit it creates in the target repo if you provide it
with authentication information in the last parameter. `AuthMethod` has a number of implementations
provided by [go-git/go-git](https://github.com/go-git/go-git), for example a `BasicAuth` structure
that lets you specify a username and password.  

### Work with an instance (thread safe)

This is the thread safe interface.

_Remember that you must pick one of the go-autumn-logging choices (see below), or call `SetupNoLoggerForTesting()` 
(not recommended) before calling into this library._

```
import (
	"context"
	generatorgit "github.com/StephanHCB/go-generator-git"
	generatorgitapi "github.com/StephanHCB/go-generator-git/api"
)

func demoCloneRenderCommit(ctx context.Context, gen generatorgitapi.GitApi) error {
	sourceUrl := "https://github.com/StephanHCB/tpl-go-rest-chi"
	sourceBranch := "master"
	targetUrl := "https://github.com/StephanHCB/scratch"
	targetBranch := "demo"
	targetBranchFrom := "main"
	generatorName := "main"
	renderSpecFile := "generated-main.yaml"
	parameters := map[string]string{} // all parameters have defaults for this generator

	if err := gen.CloneSourceRepo(ctx, sourceUrl, sourceBranch); err != nil {
		return err
	}

	if err := gen.CloneTargetRepo(ctx, targetUrl, targetBranch, targetBranchFrom); err != nil {
		return err
	}

	if _, err := gen.WriteRenderSpecFile(ctx, generatorName, renderSpecFile, parameters); err != nil {
		// find details for individual errors in first return value
		return err
	}

	if _, err := gen.Generate(ctx); err != nil {
		// find details for individual errors and the list of files that were rendered in first return value
		return err
	}

	// if auth is nil, commit won't be pushed
	if err := gen.CommitAndPush(ctx, "John Smith", "example@mailinator.com", "commit message", nil); err != nil {
		return err
	}

	return nil
}

func demoToplevel() error {
	ctx := context.TODO() // or provided from elsewhere
	basePath := "/tmp"

	gen := generatorgit.ThreadsafeInstance()

	if err := gen.CreateTemporaryWorkdir(ctx, basePath); err != nil {
		return err
	}

	if err := demoCloneRenderCommit(ctx, gen); err != nil {
		// always call Cleanup even if an error occurred to clean up after yourself
		_ = gen.Cleanup(ctx)
		return err
	}

	return gen.Cleanup(ctx)
}
```

## Implementation Prerequisites

### Choose a Logging Framework Plugin

This library uses [go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging)
to allow you to plug in the logging framework of your choice. You will need to include one of
the available specific wrappers among your dependencies. 

The simplest one, just using golang's standard
logger, is [go-autumn-logging-log](https://github.com/StephanHCB/go-autumn-logging-log).
We also have [go-autumn-logging-zerolog](https://github.com/StephanHCB/go-autumn-logging-zerolog).

If you do not want any logging, just call `aulogging.SetupNoLoggerForTesting` before calling any of the library 
functions. This will disable all logging, which is not really recommended:

```
import "github.com/StephanHCB/go-autumn-logging"

func init() {
    aulogging.SetupNoLoggerForTesting()
}
```
 
Or you can provide your own implementation of `auloggingapi.LoggingImplementation` and assign it to
`aulogging.Logger`.

## Acceptance Tests (give you examples)

We have BDD-style 
[acceptance tests](https://github.com/StephanHCB/go-generator-git/tree/master/test/acceptance). 

Running the tests and reading their code will give you lots of easy to understand examples, 
including most common error situations.
