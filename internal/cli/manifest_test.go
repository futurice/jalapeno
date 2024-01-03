package cli_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

func AddManifestSteps(s *godog.ScenarioContext) {
	s.Step(`^a manifest file that includes recipes "([^"]*)" and "([^"]*)"$`, aManifestFileThatIncludesRecipesAnd)
	s.Step(`^I execute the manifest file$`, iExecuteTheManifestFile)
}

func aManifestFileThatIncludesRecipesAnd(ctx context.Context, recipe1, recipe2 string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	dir, err := os.MkdirTemp("", "jalapeno-test-manifest")
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, manifestDirectoryPathCtxKey{}, dir)
	manifest := `apiVersion: v1
recipes:
  - foo@v0.0.1
  - bar@v0.0.1
`
	if err := os.WriteFile(filepath.Join(dir, "manifest"+re.YAMLExtension), []byte(manifest), 0644); err != nil {
		return ctx, err
	}
	return ctx, godog.ErrPending
}

func iExecuteTheManifestFile() error {
	return godog.ErrPending
}
