package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
)

func iRunPush(ctx context.Context, recipeName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"push",
		filepath.Join(recipesDir, recipeName),
		filepath.Join(ociRegistry.Resource.GetHostPort("5000/tcp"), recipeName),
	}

	if ociRegistry.TLSEnabled {
		args = append(args, "--insecure=true")
	} else {
		args = append(args, "--plain-http=true")
	}

	if ociRegistry.AuthEnabled {
		args = append(args,
			"--username=foo",
			"--password=bar",
		)
	}

	if configFileExists {
		args = append(args,
			fmt.Sprintf("--registry-config=%s", filepath.Join(configDir, DOCKER_CONFIG_FILENAME)),
		)
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}
