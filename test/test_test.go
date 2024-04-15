package cli_test

import (
	"context"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddTestSteps(s *godog.ScenarioContext) {
	s.Step(`^I run tests for recipe "([^"]*)"$`, iRunTest)
	s.Step(`^I update tests snapshosts for recipe "([^"]*)"$`, iUpdateTestSnapshot)
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
