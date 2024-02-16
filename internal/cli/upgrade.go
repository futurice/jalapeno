package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/internal/cliutil"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/conflict"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type upgradeOptions struct {
	RecipeURL      string
	ReuseOldValues bool
	SauceID        string

	option.Common
	option.OCIRepository
	option.Values
	option.WorkingDirectory
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

# Set values for new variables introduced in the upgrade
jalapeno upgrade path/to/recipe --set NEW_VAR=foo`,
	}

	cmd.Flags().StringVar(&opts.SauceID, "sauce-id", "", "If the project contains multiple sauces with the same recipe, specify the ID of the sauce to be upgraded")
	cmd.Flags().BoolVar(&opts.ReuseOldValues, "reuse-old-values", true, "Automatically set values for variables which already have a value in the existing sauce")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runUpgrade(cmd *cobra.Command, opts upgradeOptions) error {
	var (
		re       *recipe.Recipe
		oldSauce *recipe.Sauce
		err      error
	)

	if strings.HasPrefix(opts.RecipeURL, "oci://") {
		ctx := context.Background()
		re, err = recipe.PullRecipe(ctx, opts.Repository(opts.RecipeURL))

	} else {
		re, err = recipe.LoadRecipe(opts.RecipeURL)
	}

	if err != nil {
		return err
	}

	if opts.SauceID != "" {
		id, uuidErr := uuid.FromString(opts.SauceID)
		if uuidErr != nil {
			return fmt.Errorf("invalid sauce ID: %w", err)
		}

		oldSauce, err = recipe.LoadSauceByID(opts.Dir, id)
	} else {
		oldSauce, err = recipe.LoadSauceByRecipe(opts.Dir, re.Name)
		if err != nil && errors.Is(err, recipe.ErrAmbiguousSauce) {
			return fmt.Errorf(`project '%s' contains multiple sauces with recipe '%s'. Use --sauce-id to specify the ID of the sauce to be upgraded. You can check the IDs from %s/%s.%s`, opts.Dir, re.Name, recipe.SauceDirName, recipe.SaucesFileName, recipe.YAMLExtension)
		}
	}

	if err != nil {
		if errors.Is(err, recipe.ErrSauceNotFound) {
			return fmt.Errorf("project '%s' does not contain sauce with recipe '%s'. Recipe name used in the project should match the recipe which is used for upgrading", opts.Dir, re.Name)
		}

		return err
	}

	versionComparison := semver.Compare(re.Metadata.Version, oldSauce.Recipe.Metadata.Version)
	if versionComparison < 0 {
		return errors.New("new recipe version is less than the existing one")
	}

	if versionComparison > 0 {
		cmd.Printf(
			"Upgrading recipe '%s' from version %s to %s\n",
			oldSauce.Recipe.Name,
			oldSauce.Recipe.Metadata.Version,
			re.Metadata.Version,
		)
	} else {
		cmd.Printf(
			"Modifying values for sauce with recipe '%s' version %s\n",
			oldSauce.Recipe.Name,
			re.Metadata.Version,
		)
	}

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
	if opts.ReuseOtherSauceValues {
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
		cmd.Println()
		return fmt.Errorf("failed to parse provided values: %w", err)
	}

	predefinedValues := recipeutil.MergeValues(reusedValues, providedValues)

	if opts.ReuseOldValues {
		predefinedValues = recipeutil.MergeValues(oldSauce.Values, predefinedValues)
	}

	// Don't prompt variables which already has a value in existing sauce or is predefined
	varsWithoutValues := make([]recipe.Variable, 0, len(re.Variables))
	for _, v := range re.Variables {
		if _, predefinedValueExists := predefinedValues[v.Name]; !predefinedValueExists {
			varsWithoutValues = append(varsWithoutValues, v)
		}
	}

	values := predefinedValues
	if len(varsWithoutValues) > 0 {
		// If --no-input flag is set, try to use default values
		if opts.NoInput {
			varsWithoutDefaultValues := make([]recipe.Variable, 0, len(varsWithoutValues))
			for _, v := range varsWithoutValues {
				if v.Default != "" {
					defaultValue, err := v.ParseDefaultValue()
					if err != nil {
						cmd.Println()
						return fmt.Errorf("failed to parse default value for variable '%s': %w", v.Name, err)
					}
					values[v.Name] = defaultValue
				} else {
					varsWithoutDefaultValues = append(varsWithoutDefaultValues, v)
				}
			}

			// If there are still variables without values, return error
			if len(varsWithoutDefaultValues) > 0 {
				cmd.Println()
				return recipeutil.NewNoInputError(varsWithoutDefaultValues)
			}
		}

		promptedValues, err := survey.PromptUserForValues(cmd.InOrStdin(), cmd.OutOrStdout(), varsWithoutValues, predefinedValues)
		if err != nil {
			if errors.Is(err, uiutil.ErrUserAborted) {
				return nil
			}

			return fmt.Errorf("error when prompting for values: %w", err)
		}

		values = recipeutil.MergeValues(predefinedValues, promptedValues)
	}

	newSauce, err := re.Execute(engine.New(), values, oldSauce.ID)
	if err != nil {
		return err
	}

	if oldSauce.CheckFrom != "" {
		newSauce.CheckFrom = oldSauce.CheckFrom
	} else if strings.HasPrefix(opts.RecipeURL, "oci://") {
		newSauce.CheckFrom = strings.TrimSuffix(opts.RecipeURL, fmt.Sprintf(":%s", re.Metadata.Version))
	}

	newSauce.SubPath = oldSauce.SubPath

	// read common ignore file if it exists
	ignorePatterns := make([]string, 0)
	if data, err := os.ReadFile(filepath.Join(opts.Dir, recipe.IgnoreFileName)); err == nil {
		ignorePatterns = append(ignorePatterns, strings.Split(string(data), "\n")...)
	} else if !errors.Is(err, fs.ErrNotExist) {
		cmd.Println()
		// something else happened than trying to read an ignore file that does not exist
		return fmt.Errorf("failed to read ignore file: %w\n\n%s", err, cliutil.MakeRetryMessage(os.Args, values))
	}
	ignorePatterns = append(ignorePatterns, re.IgnorePatterns...)

	// Collect files which should be written to the destination directory
	output := make(map[string]recipe.File, len(newSauce.Files))
	overrideNoticed := false
	fileStatuses := make(map[string]recipeutil.FileStatus, len(newSauce.Files))

	for path := range newSauce.Files {
		skip := false
		for _, pattern := range ignorePatterns {
			if matched, err := filepath.Match(pattern, path); err != nil {
				cmd.Println()
				return fmt.Errorf("bad ignore pattern '%s': %w\n\n%s", pattern, err, cliutil.MakeRetryMessage(os.Args, values))
			} else if matched {
				// file was marked as ignored for upgrades
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		newFile := newSauce.Files[path]

		// Check if the file from previous recipe version has been modified manually
		prevFile, exists := oldSauce.Files[path]
		if exists {
			if !prevFile.HasBeenModifiedByUser() {
				output[path] = newFile
				if prevFile.Checksum == newFile.Checksum {
					fileStatuses[path] = recipeutil.FileUnchanged
				} else {
					fileStatuses[path] = recipeutil.FileModified
				}
				continue
			}

			// Check if the file has been already created manually by the user
		} else {
			existingFileContent, err := os.ReadFile(filepath.Join(opts.Dir, path))
			if errors.Is(err, os.ErrNotExist) {
				fileStatuses[path] = recipeutil.FileAdded
				output[path] = newFile
				continue
			} else if err != nil {
				return err
			}

			prevFile = recipe.NewFile(existingFileContent)
		}

		if opts.NoInput {
			return recipeutil.NewNoInputError(varsWithoutValues)
		}

		if !overrideNoticed {
			cmd.Println("\nSome of the files has been manually modified. Do you want to override the following files:")
			overrideNoticed = true
		}

		conflictResult, err := conflict.Solve(
			cmd.InOrStdin(),
			cmd.OutOrStdout(),
			path,
			prevFile.Content,
			newFile.Content,
		)

		if err != nil {
			if errors.Is(err, uiutil.ErrUserAborted) {
				cmd.Printf("User aborted\n\n%s\n", cliutil.MakeRetryMessage(os.Args, values))
				return nil
			}

			cmd.Println()
			return fmt.Errorf("error when prompting for question: %w\n\n%s", err, cliutil.MakeRetryMessage(os.Args, values))
		}

		if bytes.Equal(conflictResult, prevFile.Content) {
			fileStatuses[path] = recipeutil.FileUnchanged
			continue
		} else {
			fileStatuses[path] = recipeutil.FileModified
		}

		// NOTE: We need to save the checksum of the original file from the new sauce
		// so we would detect again if the file has been modified manually
		// when upgrading again
		newFile.Content = conflictResult

		output[path] = newFile
	}

	cmd.Println()

	newSauce.Files = output

	for filename := range oldSauce.Files {
		if _, found := newSauce.Files[filename]; !found {
			err := os.Remove(filepath.Join(opts.Dir, filename))
			if err != nil {
				return fmt.Errorf("failed to remove deprecated file '%s': %w", filename, err)
			}

			fileStatuses[filename] = recipeutil.FileDeleted
		}
	}

	err = newSauce.Save(opts.Dir)
	if err != nil {
		return err
	}

	changesFound := false
	for _, status := range fileStatuses {
		if status != recipeutil.FileUnchanged {
			changesFound = true
			cmd.Printf("Recipe upgraded %s\n", opts.Colors.Green.Render("successfully!"))
			tree := recipeutil.CreateFileTree(opts.Dir, fileStatuses)
			cmd.Printf("The following files have been processed by the recipe:\n\n%s", tree)
			break
		}
	}

	if !changesFound {
		cmd.Println("Upgrade completed, but no changes were made to any files.")
	}

	return nil
}
