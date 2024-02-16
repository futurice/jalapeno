package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func AddUpgradeSteps(s *godog.ScenarioContext) {
	s.Step(`^I upgrade recipe "([^"]*)"$`, iRunUpgrade)
	s.Step(`^I upgrade recipe from the local OCI repository "([^"]*)"$`, iRunUpgradeFromRemoteRecipe)
	s.Step(`^no conflicts were reported$`, noConflictsWereReported)
	s.Step(`^conflicts are reported$`, conflictsAreReported)
	s.Step(`^I change project file "([^"]*)" to contain "([^"]*)"$`, iChangeProjectFileToContain)
	s.Step(`^I change recipe "([^"]*)" template "([^"]*)" to render "([^"]*)"$`, iChangeRecipeTemplateToRender)
	s.Step(`^I change recipe "([^"]*)" to version "([^"]*)"$`, iChangeRecipeToVersion)
	s.Step(`^I select sauce in index (\d+) for the upgrade$`, iSelectSauceForUpgrade)
}

func iRunUpgrade(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	stdIn := ctx.Value(cmdStdInCtxKey{}).(*BlockBuffer)

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

	if stdIn.Len() == 0 {
		args = append(args, "--no-input")
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()

	ctx = clearAdditionalFlags(ctx)
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

	return iRunUpgrade(ctx, url)
}

func conflictsAreReported(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "manually modified")
}

func noConflictsWereReported(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(*bytes.Buffer)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)
	if matched, _ := regexp.MatchString("manually modified", cmdStdOut.String()); matched {
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
	recipeDir := filepath.Join(recipesDir, recipeName)

	re, err := recipe.LoadRecipe(recipeDir)
	if err != nil {
		return ctx, err
	}

	re.Version = version

	if err = re.Save(recipeDir); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func iSelectSauceForUpgrade(ctx context.Context, sauceIndex string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	sauces, err := recipe.LoadSauces(projectDir)
	if err != nil {
		return ctx, err
	}

	i, err := strconv.Atoi(sauceIndex)
	if err != nil {
		return ctx, err
	}

	additionalFlags["sauce-id"] = sauces[i].ID.String()

	return ctx, nil
}
