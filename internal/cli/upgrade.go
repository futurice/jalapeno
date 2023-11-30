package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/survey"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type upgradeOptions struct {
	RecipeURL string
	option.OCIRepository
	option.WorkingDirectory
	option.Values
	option.Common
}

func NewUpgradeCmd() *cobra.Command {
	var opts upgradeOptions
	var cmd = &cobra.Command{
		Use:   "upgrade RECIPE_URL",
		Short: "Upgrade a recipe in a project",
		Long:  "Upgrade a recipe in a project with a newer version. Recipe URL can be a local path or a remote URL (ex. 'oci://docker.io/my-recipe').",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeURL = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(cmd, opts)
		},
		Example: `# Upgrade recipe with local recipe
jalapeno upgrade path/to/recipe

# Upgrade recipe with remote recipe from OCI repository
jalapeno upgrade oci://ghcr.io/user/my-recipe:v2.0.0

# Upgrade recipe with remote recipe from OCI repository with inline authentication
jalapeno upgrade oci://ghcr.io/user/my-recipe:v2.0.0 --username user --password pass

# Upgrade recipe with remote recipe from OCI repository with Docker authentication
docker login ghcr.io
jalapeno upgrade oci://ghcr.io/user/my-recipe:v2.0.0

# Upgrade recipe to different directory
jalapeno upgrade path/to/recipe --dir other/dir

# Predefine values for new variables introduced in the upgrade
jalapeno upgrade path/to/recipe --set NEW_VAR=foo`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runUpgrade(cmd *cobra.Command, opts upgradeOptions) error {
	var (
		re  *recipe.Recipe
		err error
	)

	if strings.HasPrefix(opts.RecipeURL, "oci://") {
		ctx := context.Background()
		re, err = oci.PullRecipe(ctx, opts.Repository(opts.RecipeURL))

	} else {
		re, err = recipe.LoadRecipe(opts.RecipeURL)
	}

	if err != nil {
		return err
	}

	oldSauce, err := recipe.LoadSauce(opts.Dir, re.Name)
	if err != nil {
		if errors.Is(err, recipe.ErrSauceNotFound) {
			return fmt.Errorf("project '%s' does not contain sauce with recipe '%s'. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, re.Name)
		}

		return err
	}

	if semver.Compare(re.Metadata.Version, oldSauce.Recipe.Metadata.Version) <= 0 {
		return errors.New("new recipe version is lower or same than the existing one")
	}

	cmd.Printf(
		"Upgrading recipe %s from version %s to %s\n",
		oldSauce.Recipe.Name,
		oldSauce.Recipe.Metadata.Version,
		re.Metadata.Version,
	)

	// Check if the new version of the recipe has removed some variables
	// which existed on previous version
	for valueName := range oldSauce.Values {
		found := false
		for _, variable := range re.Variables {
			if variable.Name == valueName {
				found = true
				break
			}
		}
		if !found {
			delete(oldSauce.Values, valueName)
		}
	}

	reusedValues := make(recipe.VariableValues)
	if opts.ReuseSauceValues {
		sauces, err := recipe.LoadSauces(opts.Dir)
		if err != nil {
			return err
		}
		for _, sauce := range sauces {
			// Skip if the sauce is the one which is being upgraded
			if sauce.Recipe.Name == re.Name {
				continue
			}

			overlappingSauceValues := make(recipe.VariableValues)
			for _, v := range re.Variables {
				if val, found := sauce.Values[v.Name]; found {
					overlappingSauceValues[v.Name] = val
				}
			}

			if len(overlappingSauceValues) > 0 {
				reusedValues = recipeutil.MergeValues(reusedValues, overlappingSauceValues)
			}
		}
	}

	providedValues, err := recipeutil.ParseProvidedValues(re.Variables, opts.Values.Flags, opts.CSVDelimiter)
	if err != nil {
		return fmt.Errorf("failed to parse provided values: %w", err)
	}

	predefinedValues := recipeutil.MergeValues(reusedValues, providedValues)
	values := recipeutil.MergeValues(oldSauce.Values, predefinedValues)

	// Don't prompt variables which already has a value in existing sauce or is predefined
	varsWithoutValues := make([]recipe.Variable, 0, len(re.Variables))
	for _, v := range re.Variables {
		_, oldValueExists := oldSauce.Values[v.Name]
		_, predefinedValueExists := predefinedValues[v.Name]
		if !oldValueExists && !predefinedValueExists {
			varsWithoutValues = append(varsWithoutValues, v)
		}
	}

	if len(varsWithoutValues) > 0 {
		if opts.NoInput {
			var errMsg string
			if len(varsWithoutValues) == 1 {
				errMsg = fmt.Sprintf("value for variable %s is", varsWithoutValues[0].Name)
			} else {
				vars := make([]string, len(varsWithoutValues))
				for i, v := range varsWithoutValues {
					vars[i] = v.Name
				}
				errMsg = fmt.Sprintf("values for variables [%s] are", strings.Join(vars, ","))
			}

			return fmt.Errorf("%s missing and `--no-input` flag was set to true", errMsg)
		}

		promptedValues, err := survey.PromptUserForValues(cmd.InOrStdin(), cmd.OutOrStdout(), varsWithoutValues, predefinedValues)
		if err != nil {
			if errors.Is(err, survey.ErrUserAborted) {
				return nil
			}

			return fmt.Errorf("error when prompting for values: %w", err)
		}

		values = recipeutil.MergeValues(values, promptedValues)
	}

	newSauce, err := re.Execute(engine.Engine{}, values, oldSauce.ID)
	if err != nil {
		return err
	}

	// read common ignore file if it exists
	ignorePatterns := make([]string, 0)
	if data, err := os.ReadFile(filepath.Join(opts.Dir, recipe.IgnoreFileName)); err == nil {
		ignorePatterns = append(ignorePatterns, strings.Split(string(data), "\n")...)
	} else if !errors.Is(err, fs.ErrNotExist) {
		// something else happened than trying to read an ignore file that does not exist
		return fmt.Errorf("failed to read ignore file: %w", err)
	}
	ignorePatterns = append(ignorePatterns, re.IgnorePatterns...)

	// Collect files which should be written to the destination directory
	output := make(map[string]recipe.File, len(newSauce.Files))
	overrideNoticed := false

	for path := range newSauce.Files {
		skip := false
		for _, pattern := range ignorePatterns {
			if matched, err := filepath.Match(pattern, path); err != nil {
				return fmt.Errorf("bad ignore pattern '%s': %w", pattern, err)
			} else if matched {
				// file was marked as ignored for upgrades
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if prevFile, exists := oldSauce.Files[path]; exists {
			// Check if file was modified after rendering
			filePath := filepath.Join(opts.Dir, path)
			if modified, err := recipeutil.IsFileModified(filePath, prevFile); err != nil {
				return err
			} else if modified {
				// The file contents has been modified
				if !overrideNoticed {
					cmd.Println("Some of the files has been manually modified. Do you want to override the following files:")
					overrideNoticed = true
				}

				// TODO: We could do better in terms of merge conflict management. Like show the diff or something
				var override bool
				if err != nil {
					return fmt.Errorf("error when prompting for question: %w", err)
				}

				if !override {
					// User decided not to override the file with manual changes, remove from
					// list of changes to write
					continue
				}
			}
		}

		// Add new file or replace existing one
		output[path] = newSauce.Files[path]
	}

	newSauce.Files = output

	err = newSauce.Save(opts.Dir)
	if err != nil {
		return err
	}

	return nil
}
