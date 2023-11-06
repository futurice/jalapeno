package cli_test

import (
	"context"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunWhy(ctx context.Context, file string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewWhyCmd)

	cmd.SetArgs([]string{file})

	flags := cmd.Flags()
	if err := flags.Set("dir", projectDir); err != nil {
		return ctx, err
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			if err := flags.Set(name, value); err != nil {
				return ctx, err
			}
		}
	}

	cmd.Execute()
	return ctx, nil
}
