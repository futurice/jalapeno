package cli_test

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunWhy(ctx context.Context, file string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewWhyCmd)

	args := []string{
		file,
		fmt.Sprintf("--dir=%s", projectDir),
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		}
	}

	cmd.SetArgs(args)
	cmd.Execute()
	return ctx, nil
}
