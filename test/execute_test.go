package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func AddExecuteSteps(s *godog.ScenarioContext) {
	s.Step(`^I execute recipe "([^"]*)"$`, iRunExecute)
	s.Step(`^I execute the recipe from the local OCI repository "([^"]*)"$`, iExecuteRemoteRecipe)
	s.Step(`^recipes will be executed to the subpath "([^"]*)"$`, recipesWillBeExecutedToTheSubPath)
	s.Step(`^execution of the recipe has succeeded$`, executionOfTheRecipeHasSucceeded)
	s.Step(`^a manifest file$`, aManifestFile)
	s.Step(`^a manifest file that includes recipes$`, aManifestFileThatIncludesRecipes)
	s.Step(`^a manifest file that includes remote recipes$`, aManifestFileThatIncludesRemoteRecipes)
	s.Step(`^I execute the manifest file$`, iExecuteTheManifestFile)
	s.Step(`^I execute the manifest file with remote recipes$`, iExecuteTheManifestFileWithRemoteRecipes)
}

var TestManifestFileName = fmt.Sprintf("%s.%s", recipe.ManifestFileName, recipe.YAMLExtension)

func iRunExecute(ctx context.Context, target string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	stdIn := ctx.Value(cmdStdInCtxKey{}).(*BlockBuffer)

	var url string
	if strings.HasPrefix(target, "oci://") {
		url = target
	} else if target == TestManifestFileName {
		manifestDir := ctx.Value(manifestDirectoryPathCtxKey{}).(string)
		url = filepath.Join(manifestDir, TestManifestFileName)
	} else {
		recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
		url = filepath.Join(recipesDir, target)
	}

	args := []string{
		"execute",
		url,
		fmt.Sprintf("--dir=%s", projectDir),
	}

	if stdIn.Len() == 0 {
		args = append(args, "--no-input")
	}

	return executeCLI(ctx, args...)
}

func iExecuteRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	url := fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), repository)

	ctx = addRegistryRelatedFlags(ctx)

	return iRunExecute(ctx, url)
}

func recipesWillBeExecutedToTheSubPath(ctx context.Context, path string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["subpath"] = path

	return ctx, nil
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}

func aManifestFile(ctx context.Context) (context.Context, error) {
	manifestDir := ctx.Value(manifestDirectoryPathCtxKey{}).(string)
	manifest := recipe.NewManifest()

	err := manifest.Save(filepath.Join(manifestDir, TestManifestFileName))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func aManifestFileThatIncludesRecipes(ctx context.Context, recipeNames *godog.Table) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	manifestDir := ctx.Value(manifestDirectoryPathCtxKey{}).(string)

	recipes := make([]recipe.ManifestRecipe, len(recipeNames.Rows))
	for i, row := range recipeNames.Rows {
		name := row.Cells[0].Value
		re, err := recipe.LoadRecipe(filepath.Join(recipeDir, name))
		if err != nil {
			return ctx, err
		}

		recipes[i] = recipe.ManifestRecipe{
			Name:       re.Name,
			Version:    re.Version,
			Repository: fmt.Sprintf("file://%s/%s", recipeDir, name),
		}
	}

	manifest := recipe.NewManifest()
	manifest.Recipes = recipes

	err := manifest.Save(filepath.Join(manifestDir, TestManifestFileName))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func aManifestFileThatIncludesRemoteRecipes(ctx context.Context, recipeNames *godog.Table) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	manifestDir := ctx.Value(manifestDirectoryPathCtxKey{}).(string)

	recipes := make([]recipe.ManifestRecipe, len(recipeNames.Rows))
	for i, row := range recipeNames.Rows {
		name := row.Cells[0].Value
		version := row.Cells[1].Value

		recipes[i] = recipe.ManifestRecipe{
			Name:       name,
			Version:    version,
			Repository: fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), name),
		}
	}

	manifest := recipe.NewManifest()
	manifest.Recipes = recipes

	err := manifest.Save(filepath.Join(manifestDir, TestManifestFileName))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func iExecuteTheManifestFile(ctx context.Context) (context.Context, error) {
	return iRunExecute(ctx, TestManifestFileName)
}

func iExecuteTheManifestFileWithRemoteRecipes(ctx context.Context) (context.Context, error) {
	ctx = addRegistryRelatedFlags(ctx)
	return iRunExecute(ctx, TestManifestFileName)
}
