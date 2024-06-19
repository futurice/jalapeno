package cli_test

import (
	"context"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func AddTestSteps(s *godog.ScenarioContext) {
	s.Step(`^I run tests for recipe "([^"]*)"$`, iRunTest)
	s.Step(`^I update tests snapshosts for recipe "([^"]*)"$`, iUpdateTestSnapshot)
	s.Step(`^I expect recipe "([^"]*)" initHelp to match "([^"]*)"$`, iExpectRecipesInitHelpToMatch)
}

func iRunTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"test",
		filepath.Join(recipesDir, recipe),
	)
}

func iUpdateTestSnapshot(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["update-snapshots"] = "true"

	return iRunTest(ctx, recipe)
}

func iExpectRecipesInitHelpToMatch(ctx context.Context, recipeName, expectedInitHelp string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipeDir, recipeName))
	if err != nil {
		return ctx, err
	}

	re.Tests[0].ExpectedInitHelp = expectedInitHelp
	if err := re.Save(filepath.Join(recipeDir, recipeName)); err != nil {
		return ctx, err
	}

	return ctx, nil
}
