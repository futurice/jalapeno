package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddPullSteps(s *godog.ScenarioContext) {
	s.Step(`^I pull the recipe "([^"]*)" from the local OCI repository "([^"]*)"$`, iPullRecipe)
}

func iPullRecipe(ctx context.Context, recipeName, repoName string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)

	args := []string{
		"pull",
		filepath.Join(ociRegistry.Resource.GetHostPort("5000/tcp"), repoName),
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
