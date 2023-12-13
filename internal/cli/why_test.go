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
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"why",
		file,
		fmt.Sprintf("--dir=%s", projectDir),
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}
