package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func iRunUpgrade(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	var url string
	if strings.HasPrefix(recipe, "oci://") {
		url = recipe
	} else {
		recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
		url = filepath.Join(recipesDir, recipe)
	}

	args := []string{
		"upgrade",
		url,
		fmt.Sprintf("--dir=%s", projectDir),
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func iRunUpgradeFromRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	url := fmt.Sprintf("oci://%s/%s", registry.Resource.GetHostPort("5000/tcp"), repository)

	if registry.TLSEnabled {
		// Allow self-signed certificates
		additionalFlags["insecure"] = "true"
	} else {
		additionalFlags["plain-http"] = "true"
	}

	if registry.AuthEnabled {
		additionalFlags["username"] = "foo"
		additionalFlags["password"] = "bar"
	}

	if configFileExists && os.Getenv("DOCKER_CONFIG") == "" {
		additionalFlags["registry-config"] = filepath.Join(configDir, DOCKER_CONFIG_FILENAME)
	}

	ctx = context.WithValue(ctx, cmdAdditionalFlagsCtxKey{}, additionalFlags)

	return iRunUpgrade(ctx, url)
}

func conflictsAreReported(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "modified")
}

func noConflictsWereReported(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(*bytes.Buffer)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)
	if matched, _ := regexp.MatchString("modified", cmdStdOut.String()); matched {
		return ctx, fmt.Errorf("Conflict in recipe\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func iChangeProjectFileToContain(ctx context.Context, filename, content string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	if err := os.WriteFile(filepath.Join(projectDir, filename), []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iChangeRecipeTemplateToRender(ctx context.Context, recipeName, filename, content string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	templateFilePath := filepath.Join(recipesDir, recipeName, "templates", filename)
	if err := os.WriteFile(templateFilePath, []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iChangeRecipeToVersion(ctx context.Context, recipeName, version string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	re, err := recipe.LoadRecipe(filepath.Join(recipesDir, recipeName))
	if err != nil {
		return ctx, err
	}

	re.Version = version

	if err = re.Save(recipesDir); err != nil {
		return ctx, err
	}

	return ctx, nil
}
