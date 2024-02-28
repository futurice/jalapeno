package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/internal/cliutil"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type executeOptions struct {
	RecipeURL string
	Subpath   string

	option.Common
	option.OCIRepository
	option.Values
	option.WorkingDirectory
}

func NewExecuteCmd() *cobra.Command {
	var opts executeOptions
	var cmd = &cobra.Command{
		Use:     "execute RECIPE_URL",
		Aliases: []string{"exec", "e", "run"},
		Short:   "Execute a recipe",
		Long:    "Executes (renders) a recipe and outputs the files to the directory. Recipe URL can be a local path or a remote URL (ex. 'oci://docker.io/my-recipe').",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeURL = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExecute(cmd, opts)
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

	if err := recipe.ValidateSubpath(opts.Subpath); err != nil {
		return err
	}

	var (
		re              *recipe.Recipe
		err             error
		wasRemoteRecipe bool
	)

	if strings.HasPrefix(opts.RecipeURL, "oci://") {
		wasRemoteRecipe = true
		ctx := context.Background()
		re, err = recipe.PullRecipe(ctx, opts.Repository(opts.RecipeURL))

	} else {
		re, err = recipe.LoadRecipe(opts.RecipeURL)
	}

	if err != nil {
		return fmt.Errorf("can not load the recipe: %s", err)
	}

	cmd.Printf("%s: %s\n", opts.Colors.Red.Render("Recipe name"), re.Metadata.Name)

	if re.Metadata.Description != "" {
		cmd.Printf("%s: %s\n", opts.Colors.Red.Render("Description"), re.Metadata.Description)
	}

	cmd.Println()

	// Load all existing sauces
	existingSauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	for _, sauce := range existingSauces {
		if sauce.Recipe.Name == re.Name &&
			semver.Compare(sauce.Recipe.Metadata.Version, re.Metadata.Version) == 0 &&
			sauce.SubPath == opts.Subpath {
			return fmt.Errorf("recipe '%s' with version '%s' has been already executed. If you want to re-execute the recipe with different values, use `upgrade` command with `--reuse-old-values=false` flag", re.Name, re.Metadata.Version)
		}
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
			if errors.Is(err, uiutil.ErrUserAborted) {
				return nil
			} else {
				return fmt.Errorf("error when prompting for values: %s", err)
			}
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

	sauce.SubPath = opts.Subpath

	// Automatically add recipe origin if the recipe was remote
	if wasRemoteRecipe {
		sauce.CheckFrom = strings.TrimSuffix(opts.RecipeURL, fmt.Sprintf(":%s", re.Metadata.Version))
	}

	// Check for conflicts
	for _, s := range existingSauces {
		if conflicts := s.Conflicts(sauce); conflicts != nil {
			retryMessage := cliutil.MakeRetryMessage(os.Args, values)
			return fmt.Errorf("conflict in recipe '%s': file '%s' was already created by recipe '%s'.\n\n%s", re.Name, conflicts[0].Path, s.Recipe.Name, retryMessage)
		}
	}

	err = sauce.Save(opts.Dir)
	if err != nil {
		return err
	}

	cmd.Printf("Recipe executed %s\n", opts.Colors.Green.Render("successfully!"))

	files := sauce.Files
	if opts.Subpath != "" {
		files = make(map[string]recipe.File, len(sauce.Files))
		for path, file := range sauce.Files {
			files[filepath.Join(opts.Subpath, path)] = file
		}
	}

	fileTreeFiles := make(map[string]recipeutil.FileStatus, len(files))
	for path := range files {
		fileTreeFiles[path] = recipeutil.FileAdded
	}

	tree := recipeutil.CreateFileTree(opts.Dir, fileTreeFiles)
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
