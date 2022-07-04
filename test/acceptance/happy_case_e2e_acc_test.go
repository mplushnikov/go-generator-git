package acceptance

import (
	"context"
	generatorgit "github.com/StephanHCB/go-generator-git/v2"
	"github.com/StephanHCB/go-generator-git/v2/docs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHappyPath_End2End_NewTargetBranch(t *testing.T) {
	docs.Given("a valid generator source and target repository")
	sourceUrl := "https://github.com/StephanHCB/tpl-go-rest-chi"
	sourceBranch := "master"
	targetUrl := "https://github.com/StephanHCB/scratch"
	targetBranch := "test-e2e-happy-path-1-nopush" // means it'll never exist
	targetFrom := "main"
	generatorName := "main"
	renderSpecFile := "generated-main.yaml"
	parameters := map[string]interface{}{} // all parameters have defaults for this generator

	ctx := context.TODO()

	docs.When("the git generator is invoked")
	docs.Then("no errors occur")
	err := generatorgit.CreateTemporaryWorkdir(ctx, "../output")
	require.Nil(t, err)

	path, err := generatorgit.CloneSourceRepo(ctx, sourceUrl, sourceBranch, nil)
	require.NotNil(t, path)
	require.Nil(t, err)

	path, err = generatorgit.CloneTargetRepo(ctx, targetUrl, targetBranch, targetFrom, nil)
	require.NotNil(t, path)
	require.Nil(t, err)

	response, err := generatorgit.WriteRenderSpecFile(ctx, generatorName, renderSpecFile, parameters)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.True(t, response.Success)
	// TODO check response some more

	response, err = generatorgit.Generate(ctx)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.True(t, response.Success)
	// TODO check response some more, contains a certain file? No errors?

	docs.Then("the repositories are cloned as expected and rendering succeeds")
	// TODO check genspec, renderspec, and one other small file

	docs.Then("commit and (simulated) push succeed")
	err = generatorgit.CommitAndPush(ctx, "somebody", "somebody@mailinator.com", "initial generation", nil)
	require.Nil(t, err)
	// TODO check that no open changes in target repo any more
	// TODO check that new commit was made

	docs.Then("no spurious output remains")
	err = generatorgit.Cleanup(ctx)
	require.Nil(t, err)
	// TODO check that output has no subdirectories
}
