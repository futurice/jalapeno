package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func iPullRecipe(ctx context.Context, recipeName, repoName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)

	cmd, cmdStdOut, cmdStdErr := wrapCmdOutputs(newPullCmd)

	cmd.SetArgs([]string{filepath.Join(registry.Resource.GetHostPort("5000/tcp"), repoName)})
	flags := cmd.Flags()
	if err := flags.Set("output", recipesDir); err != nil {
		return ctx, err
	}
	if registry.TLSEnabled {
		// Allow self-signed certificates
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
			return ctx, err
		}
	}

	if registry.AuthEnabled {
		if err := flags.Set("username", "foo"); err != nil {
			return ctx, err
		}
		if err := flags.Set("password", "bar"); err != nil {
			return ctx, err
		}
	}

	if configFileExists {
		if err := flags.Set("registry-config", filepath.Join(configDir, DOCKER_CONFIG_FILENAME)); err != nil {
			return ctx, err
		}
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func pullOfTheRecipeWasSuccessful(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)

	if cmdStdErr != "" {
		return ctx, fmt.Errorf("stderr was not empty: %s", cmdStdErr)
	}

	if cmdStdOut == "" { // TODO: Check stdout when we have proper message from CMD
		return ctx, errors.New("stdout was empty")
	}

	return ctx, nil
}

func pullOfTheRecipeHasFailedWithError(ctx context.Context, errorMessage string) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)

	if cmdStdOut != "" {
		return ctx, fmt.Errorf("stdout was not empty: %s", cmdStdErr)
	}

	if cmdStdErr == "" {
		return ctx, errors.New("stderr was empty")
	}

	if strings.TrimSpace(cmdStdErr) != errorMessage {
		return ctx, fmt.Errorf("error message did not match: expected '%s', found '%s", errorMessage, cmdStdErr)
	}

	return ctx, nil
}
