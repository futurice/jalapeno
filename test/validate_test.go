package cli_test

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func AddValidateSteps(s *godog.ScenarioContext) {
	s.Step(`I validate recipe "([^"]*)"$`, iRunValidate)
}

func iRunValidate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"validate",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	)
}
