package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddCreateSteps(s *godog.ScenarioContext) {
	s.Step(`^I create a recipe with name "([^"]*)"$`, iRunCreateRecipe)
	s.Step(`^I create a test for recipe "([^"]*)"$`, iRunCreateTest)
}

func iRunCreateRecipe(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"create",
		"recipe",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	)
}

func iRunCreateManifest(ctx context.Context) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"create",
		"manifest",
		fmt.Sprintf("--dir=%s", recipesDir),
	)
}

func iRunCreateTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"create",
		"test",
		fmt.Sprintf("--dir=%s", filepath.Join(recipesDir, recipe)),
	)
}
