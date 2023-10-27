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
		Use:   "upgrade RECIPE_PATH",
		Short: "Upgrade a recipe in a project",
		Long:  "Upgrade a recipe in a project with a newer version.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeURL = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runUpgrade(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runUpgrade(cmd *cobra.Command, opts upgradeOptions) {
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
		cmd.PrintErrf("Error: %s", err)
		return
	}

	oldSauce, err := recipe.LoadSauce(opts.Dir, re.Name)
	if err != nil {
		var msg string
		if errors.Is(err, recipe.ErrSauceNotFound) {
			msg = fmt.Sprintf("project '%s' does not contain sauce with recipe '%s'. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, re.Name)
		} else {
			msg = err.Error()
		}
		cmd.PrintErrf("Error: %s", msg)
		return
	}

	if semver.Compare(re.Metadata.Version, oldSauce.Recipe.Metadata.Version) <= 0 {
		cmd.PrintErrln("Error: new recipe version is lower or same than the existing one")
		return
	}

	cmd.Printf("Upgrading recipe %s from version %s to %s\n", oldSauce.Recipe.Name, oldSauce.Recipe.Metadata.Version, re.Metadata.Version)

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
			cmd.PrintErrf("Error: %s", err)
			return
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
		cmd.PrintErrf("Error when parsing provided values: %v\n", err)
		return
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

			cmd.PrintErrf("Error: %s missing and `--no-input` flag was set to true\n", errMsg)
			return
		}

		promptedValues, err := survey.PromptUserForValues(cmd.InOrStdin(), cmd.OutOrStdout(), varsWithoutValues, predefinedValues)
		if err != nil {
			if !errors.Is(err, survey.ErrUserAborted) {
				cmd.PrintErrf("Error when prompting for values: %s\n", err)
			}
			return
		}

		values = recipeutil.MergeValues(values, promptedValues)
	}

	newSauce, err := re.Execute(values, oldSauce.ID)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	// read common ignore file if it exists
	ignorePatterns := make([]string, 0)
	if data, err := os.ReadFile(filepath.Join(opts.Dir, recipe.IgnoreFileName)); err == nil {
		ignorePatterns = append(ignorePatterns, strings.Split(string(data), "\n")...)
	} else if !errors.Is(err, fs.ErrNotExist) {
		// something else happened than trying to read an ignore file that does not exist
		cmd.PrintErrf("Error: failed to read ignore file: %s\n", err)
		return
	}
	ignorePatterns = append(ignorePatterns, re.IgnorePatterns...)

	// Collect files which should be written to the destination directory
	output := make(map[string]recipe.File, len(newSauce.Files))
	overrideNoticed := false

	for path := range newSauce.Files {
		skip := false
		for _, pattern := range ignorePatterns {
			if matched, err := filepath.Match(pattern, path); err != nil {
				cmd.PrintErrf("Error: bad ignore pattern '%s': %s\n", pattern, err)
				return
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
				cmd.PrintErrf("Error: %s", err)
				return
			} else if modified {
				// The file contents has been modified
				if !overrideNoticed {
					cmd.Println("Some of the files has been manually modified. Do you want to override the following files:")
					overrideNoticed = true
				}

				// TODO: We could do better in terms of merge conflict management. Like show the diff or something
				var override bool
				// prompt := &survey.Confirm{
				// 	Message: path,
				// 	Default: true,
				// }

				// err = survey.AskOne(prompt, &override)
				if err != nil {
					cmd.PrintErrf("Error when prompting for question: %s", err)
					return
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
		cmd.PrintErrf("Error: %s", err)
		return
	}
}
