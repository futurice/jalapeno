package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddTestSteps(s *godog.ScenarioContext) {
	s.Step(`^I run tests for recipe "([^"]*)"$`, iRunTest)
	s.Step(`^I run tests for recipe "([^"]*)" while updating snapshots$`, iRunTestWithSnapshotUpdate)
	s.Step(`^I create a test for recipe "([^"]*)"$`, iCreateRecipeTest)
}

func iRunTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"test",
		filepath.Join(recipesDir, recipe),
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func iCreateRecipeTest(ctx context.Context, recipe string) (context.Context, error) {
	ctx = context.WithValue(
		ctx,
		cmdAdditionalFlagsCtxKey{},
		map[string]string{"create": "true"},
	)

	return iRunTest(ctx, recipe)
}

func iRunTestWithSnapshotUpdate(ctx context.Context, recipe string) (context.Context, error) {
	ctx = context.WithValue(
		ctx,
		cmdAdditionalFlagsCtxKey{},
		map[string]string{"update-snapshots": "true"},
	)

	return iRunTest(ctx, recipe)
}
