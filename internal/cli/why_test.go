package cli_test

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func AddWhySteps(s *godog.ScenarioContext) {
	s.Step(`I check why the file "([^"]*)" is created$`, iRunWhy)
}

func iRunWhy(ctx context.Context, file string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"why",
		file,
		fmt.Sprintf("--dir=%s", projectDir),
	)
}
