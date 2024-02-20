package cli_test

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func AddCreateSteps(s *godog.ScenarioContext) {
	s.Step(`^I create a recipe with name "([^"]*)"$`, iRunCreate)
}

func iRunCreate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"create",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	)
}
