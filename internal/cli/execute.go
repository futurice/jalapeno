package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/survey"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

type executeOptions struct {
	RecipeURL string
	option.Values
	option.OCIRepository
	option.WorkingDirectory
	option.Common
}

func NewExecuteCmd() *cobra.Command {
	var opts executeOptions
	var cmd = &cobra.Command{
		Use:     "execute RECIPE_URL",
		Aliases: []string{"exec", "e"},
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

# Predefine variable values
jalapeno execute path/to/recipe --set MY_VAR=foo --set MY_OTHER_VAR=bar`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runExecute(cmd *cobra.Command, opts executeOptions) error {
	if _, err := os.Stat(opts.Dir); os.IsNotExist(err) {
		return errors.New("output path does not exist")
	}

	var (
		re              *recipe.Recipe
		err             error
		wasRemoteRecipe bool
	)

	if strings.HasPrefix(opts.RecipeURL, "oci://") {
		wasRemoteRecipe = true
		ctx := context.Background()
		re, err = oci.PullRecipe(ctx, opts.Repository(opts.RecipeURL))

	} else {
		re, err = recipe.LoadRecipe(opts.RecipeURL)
	}

	if err != nil {
		return fmt.Errorf("can not load the recipe: %s", err)
	}

	style := lipgloss.NewStyle().Foreground(opts.Colors.Primary)
	cmd.Printf("%s: %s\n", style.Render("Recipe name"), re.Metadata.Name)

	if re.Metadata.Description != "" {
		cmd.Printf("%s: %s\n", style.Render("Description"), re.Metadata.Description)
	}

	// Load all existing sauces
	existingSauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	reusedValues := make(recipe.VariableValues)
	if opts.ReuseSauceValues && len(existingSauces) > 0 {
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

	providedValues, err := recipeutil.ParseProvidedValues(re.Variables, opts.Values.Flags, opts.Values.CSVDelimiter)
	if err != nil {
		return fmt.Errorf("failed to parse provided values: %w", err)
	}

	values := recipeutil.MergeValues(reusedValues, providedValues)

	// Filter out variables which don't have value yet
	varsWithoutValues := recipeutil.FilterVariablesWithoutValues(re.Variables, values)
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

		promptedValues, err := survey.PromptUserForValues(cmd.InOrStdin(), cmd.OutOrStdout(), varsWithoutValues, values)
		if err != nil {
			if errors.Is(err, survey.ErrUserAborted) {
				return nil
			} else {
				return fmt.Errorf("error when prompting for values: %s", err)
			}
		}
		values = recipeutil.MergeValues(values, promptedValues)
	}

	sauce, err := re.Execute(values, uuid.Must(uuid.NewV4()))
	if err != nil {
		return err
	}

	// Check for conflicts
	for _, s := range existingSauces {
		if conflicts := s.Conflicts(sauce); conflicts != nil {
			return fmt.Errorf("conflict in recipe '%s': file '%s' was already created by recipe '%s'", re.Name, conflicts[0].Path, s.Recipe.Name)
		}
	}

	// Automatically add recipe origin if the recipe was remote
	if wasRemoteRecipe {
		sauce.CheckFrom = opts.RecipeURL
	}

	err = sauce.Save(opts.Dir)
	if err != nil {
		return err
	}

	cmd.Println("\nRecipe executed successfully!")

	tree := recipeutil.CreateFileTree(opts.Dir, sauce.Files)
	cmd.Printf("The following files were created:\n\n%s", tree)

	if re.InitHelp != "" {
		cmd.Printf("\nNext up: %s\n", re.InitHelp)
	}

	return nil
}
