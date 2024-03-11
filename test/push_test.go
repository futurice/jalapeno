package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddPushSteps(s *godog.ScenarioContext) {
	s.Step(`^I push the recipe "([^"]*)" to the local OCI repository$`, iRunPush)
	s.Step(`^I push the recipe "([^"]*)" to the local OCI repository with \'--latest\' flag$`, iRunPushWithLatestTag)
	s.Step(`^the recipe "([^"]*)" is pushed to the local OCI repository "([^"]*)"$`, iRunPush)
}

func iRunPush(ctx context.Context, recipeName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)

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

	return executeCLI(ctx, args...)
}

func iRunPushWithLatestTag(ctx context.Context, recipeName string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["latest"] = "true"

	return iRunPush(ctx, recipeName)
}
