# go-generator-git

A golang library for generating files from templates, capable of obtaining the templates from one git repository
and committing and pushing the result to another git repository. Can be used to scaffold application code.

See [go-generator-lib](https://github.com/StephanHCB/go-generator-lib/) for details regarding the templates and
the generation process. This library just clones the git repositories to two directories, 
creates or checks out the necessary branches, and handles the final commit and push in the target repository. 
It relies on _go-generator-lib_ to generate output files between the two directories.

## Usage

```
TODO
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

```
TODO
```

We have BDD-style 
[acceptance tests](https://github.com/StephanHCB/go-generator-git/tree/master/test/acceptance). 

Running the tests and reading their code will give you lots of easy to understand examples, 
including most common error situations.
