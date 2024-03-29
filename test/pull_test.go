package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddPullSteps(s *godog.ScenarioContext) {
	s.Step(`^I pull recipe from the local OCI repository "([^"]*)"$`, iPullRecipe)
}

func iPullRecipe(ctx context.Context, repoName string) (context.Context, error) {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)

	args := []string{
		"pull",
		filepath.Join(ociRegistry.Resource.GetHostPort("5000/tcp"), repoName),
		fmt.Sprintf("--dir=%s", dir),
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
