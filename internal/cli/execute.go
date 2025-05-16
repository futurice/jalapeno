package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/internal/cliutil"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

type executeOptions struct {
	RecipeURL string
	Subpath   string

	option.Common
	option.OCIRepository
	option.Values
	option.WorkingDirectory
	option.Timeout
}

func NewExecuteCmd() *cobra.Command {
	var opts executeOptions
	var cmd = &cobra.Command{
		Use:     "execute RECIPE_URL",
		Aliases: []string{"exec", "e", "run"},
		Short:   "Execute a recipe or manifest",
		Long:    "Executes (renders) a recipe or manifest and outputs the files to the directory. Recipe URL can be a local path or a remote URL (ex. 'oci://docker.io/my-recipe').",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeURL = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runExecute(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `# Execute local recipe
jalapeno execute path/to/recipe

# Execute recipe from OCI repository
jalapeno execute oci://ghcr.io/user/my-recipe:latest

# Execute recipe from OCI repository with inline authentication
jalapeno execute oci://ghcr.io/user/my-recipe:latest --username user --password pass

# Execute recipe from OCI repository with Docker authentication
docker login ghcr.io
jalapeno execute oci://ghcr.io/user/my-recipe:latest

# Execute a manifest which contains multiple recipes
jalapeno execute path/to/manifest.yml

# Execute recipe to different directory
jalapeno execute path/to/recipe --dir other/dir

# Execute recipe in a monorepo
jalapeno execute path/to/recipe --dir monorepo-root --subpath path/in/monorepo

# Set variable values with flags
jalapeno execute path/to/recipe --set MY_VAR=foo --set MY_OTHER_VAR=bar

# Set variable values with environment variables
export JALAPENO_VAR_MY_VAR=foo
jalapeno execute path/to/recipe`,
	}

	cmd.Flags().StringVar(&opts.Subpath, "subpath", "", "Subpath which is used as a path prefix when saving and loading the sauce files. Useful for monorepos.")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runExecute(cmd *cobra.Command, opts executeOptions) error {
	if _, err := os.Stat(opts.Dir); os.IsNotExist(err) {
		return errors.New("output path does not exist")
	}

	if err := recipe.ValidateSubPath(opts.Subpath); err != nil {
		return err
	}

	var (
		re  *recipe.Recipe
		err error
	)

	switch recipe.DetermineRecipeURLType(opts.RecipeURL) {
	case recipe.OCIType:
		ctx, cancel := context.WithTimeout(cmd.Context(), opts.Timeout.Duration)
		defer cancel()

		re, err = recipe.PullRecipe(ctx, opts.Repository(opts.RecipeURL))
		if err != nil {
			return fmt.Errorf("can not load the remote recipe: %s", err)
		}
		return executeRecipe(cmd, opts, re)

	case recipe.LocalType:
		re, err = recipe.LoadRecipe(opts.RecipeURL)
		if err != nil {
			return fmt.Errorf("can not load the recipe: %s", err)
		}
		return executeRecipe(cmd, opts, re)

	case recipe.ManifestType:
		manifest, err := recipe.LoadManifest(opts.RecipeURL)
		if err != nil {
			return fmt.Errorf("can not load the manifest: %s", err)
		}
		return executeManifest(cmd, opts, manifest)

	default:
		return fmt.Errorf("unsupported recipe URL: %s", opts.RecipeURL)
	}
}

func executeRecipe(cmd *cobra.Command, opts executeOptions, re *recipe.Recipe) error {
	cmd.Printf("%s: %s\n", colors.Red.Render("Recipe name"), re.Metadata.Name)
	cmd.Printf("%s: %s\n", colors.Red.Render("Version"), re.Metadata.Version)

	if re.Metadata.Description != "" {
		cmd.Printf("%s: %s\n", colors.Red.Render("Description"), re.Metadata.Description)
	}

	cmd.Println()

	// Load all existing sauces
	existingSauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	reusedValues := make(recipe.VariableValues)
	if opts.ReuseOtherSauceValues && len(existingSauces) > 0 {
		for _, sauce := range existingSauces {
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

	providedValues, err := recipeutil.ParseProvidedValues(
		re.Variables,
		opts.Values.Flags,
		opts.Values.CSVDelimiter,
		opts.Values.ParseEnvironmentVariables,
	)
	if err != nil {
		return fmt.Errorf("failed to parse provided values: %w", err)
	}

	values := recipeutil.MergeValues(reusedValues, providedValues)

	// Filter out variables which don't have value yet
	varsWithoutValues := recipeutil.FilterVariablesWithoutValues(re.Variables, values)
	if len(varsWithoutValues) > 0 {
		// If --no-input flag is set, try to use default values
		if opts.NoInput {
			varsWithoutDefaultValues := make([]recipe.Variable, 0, len(varsWithoutValues))
			for _, v := range varsWithoutValues {
				if v.Default != "" {
					defaultValue, err := v.ParseDefaultValue()
					if err != nil {
						return fmt.Errorf("failed to parse default value for variable '%s': %w", v.Name, err)
					}
					values[v.Name] = defaultValue
				} else {
					varsWithoutDefaultValues = append(varsWithoutDefaultValues, v)
				}
			}

			// If there are still variables without values, return error
			if len(varsWithoutDefaultValues) > 0 {
				return recipeutil.NewNoInputError(varsWithoutDefaultValues)
			}
		}

		promptedValues, err := survey.PromptUserForValues(cmd.InOrStdin(), cmd.OutOrStdout(), varsWithoutValues, values)
		if err != nil {
			return fmt.Errorf("error when prompting for values: %w", err)
		}
		values = recipeutil.MergeValues(values, promptedValues)
	} else {
		cmd.Println()
	}

	sauce, err := re.Execute(engine.New(), values, uuid.Must(uuid.NewV4()))
	if err != nil {
		retryMessage := cliutil.MakeRetryMessage(os.Args, values)
		return fmt.Errorf("%w\n\n%s", err, retryMessage)
	}

	sauce.Subpath = filepath.ToSlash(opts.Subpath)

	// Automatically add recipe origin if the recipe was remote
	if recipe.DetermineRecipeURLType(opts.RecipeURL) == recipe.OCIType {
		// strip the tag from the URL
		re := regexp.MustCompile(`(.+)(:[^\/\/].+)$`)
		sauce.CheckFrom = re.ReplaceAllString(opts.RecipeURL, "$1")
	}

	// Check for conflicts
	for _, s := range existingSauces {
		if conflicts := s.Conflicts(sauce); conflicts != nil {
			retryMessage := cliutil.MakeRetryMessage(os.Args, values)
			return fmt.Errorf("conflict in recipe '%s': file '%s' was already created by other recipe '%s' (sauce ID: %s).\n\n%s", re.Name, conflicts[0].Path, s.Recipe.Name, s.ID, retryMessage)
		}
	}

	err = sauce.Save(opts.Dir)
	if err != nil {
		return err
	}

	cmd.Printf("Recipe executed %s\n", colors.Green.Render("successfully!"))

	fileTreeFiles := make(map[string]recipeutil.FileStatus, len(sauce.Files))
	for path := range sauce.Files {
		fileTreeFiles[path] = recipeutil.FileAdded
	}

	root := opts.Dir
	if opts.Subpath != "" {
		root = filepath.ToSlash(filepath.Join(root, opts.Subpath))
	}

	tree := recipeutil.CreateFileTree(root, fileTreeFiles)
	cmd.Printf("The following files were created:\n\n%s", tree)

	if re.InitHelp != "" {
		help, err := sauce.RenderInitHelp()
		if err != nil {
			return err
		}
		cmd.Printf("\nNext up: %s\n", help)
	}

	return nil
}

func executeManifest(cmd *cobra.Command, opts executeOptions, manifest *recipe.Manifest) error {
	if len(manifest.Recipes) > 1 {
		cmd.Printf("Executing manifest with %d recipes...\n\n", len(manifest.Recipes))
	}

	if len(opts.Values.Flags) > 0 {
		return errors.New("values can not be provided when executing a manifest. Use values in the manifest file instead")
	}

	for i, manifestRecipe := range manifest.Recipes {
		var re *recipe.Recipe
		var err error
		ctx, cancel := context.WithTimeout(cmd.Context(), opts.Timeout.Duration)
		defer cancel()

		switch recipe.DetermineRecipeURLType(manifestRecipe.Repository) {
		case recipe.OCIType:
			re, err = recipe.PullRecipe(
				ctx,
				opts.Repository(fmt.Sprintf("%s:%s", manifestRecipe.Repository, manifestRecipe.Version)),
			)

		case recipe.LocalType:
			re, err = recipe.LoadRecipe(manifestRecipe.Repository)
		}

		if err != nil {
			return fmt.Errorf("can not load the recipe '%s': %s", manifestRecipe.Name, err)
		}

		// Apply values provided by the manifest
		valueFlags := make([]string, 0, len(manifestRecipe.Values))
		for name, value := range manifestRecipe.Values {
			valueFlags = append(valueFlags, fmt.Sprintf("%s=%s", name, value))
		}

		opts.Values.Flags = valueFlags

		if err := executeRecipe(cmd, opts, re); err != nil {
			return err
		}

		if i < len(manifest.Recipes)-1 {
			cmd.Print("\n- - - - - - - - - -\n\n")
		}
	}

	return nil
}
