package cli_test

import (
	"context"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddTestSteps(s *godog.ScenarioContext) {
	s.Step(`^I run tests for recipe "([^"]*)"$`, iRunTest)
	s.Step(`^I update tests snapshosts for recipe "([^"]*)"$`, iUpdateTestSnapshot)
	s.Step(`^I create a test for recipe "([^"]*)"$`, iCreateRecipeTest)
}

func iRunTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"test",
		filepath.Join(recipesDir, recipe),
	)
}

func iCreateRecipeTest(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["create"] = "true"

	return iRunTest(ctx, recipe)
}

func iUpdateTestSnapshot(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["update-snapshots"] = "true"

	return iRunTest(ctx, recipe)
}
