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

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"create",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}
