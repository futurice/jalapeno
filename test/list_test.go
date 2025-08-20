package cli_test

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func AddListSteps(s *godog.ScenarioContext) {
	s.Step(`^I list all sauces in the project$`, iRunList)
}

func iRunList(ctx context.Context) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"list",
		fmt.Sprintf("--dir=%s", projectDir),
	)
}
