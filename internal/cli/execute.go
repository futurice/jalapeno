package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
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
			return recipeutil.NewNoInputError(varsWithoutValues)
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
	}

	sauce, err := re.Execute(engine.Engine{}, values, uuid.Must(uuid.NewV4()))
	if err != nil {
		retryMessage := makeRetryMessage(opts, values)
		return fmt.Errorf("%w\n\n%s", err, retryMessage)
	}

	// Check for conflicts
	for _, s := range existingSauces {
		if conflicts := s.Conflicts(sauce); conflicts != nil {
			retryMessage := makeRetryMessage(opts, values)
			return fmt.Errorf("conflict in recipe '%s': file '%s' was already created by recipe '%s'.\n\n%s", re.Name, conflicts[0].Path, s.Recipe.Name, retryMessage)
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

	cmd.Println("Recipe executed successfully!")

	tree := recipeutil.CreateFileTree(opts.Dir, sauce.Files)
	cmd.Printf("The following files were created:\n\n%s", tree)

	if re.InitHelp != "" {
		cmd.Printf("\nNext up: %s\n", re.InitHelp)
	}

	return nil
}

func makeRetryMessage(opts executeOptions, values recipe.VariableValues) string {
	var commandline strings.Builder
	commandline.WriteString("jalapeno execute ")
	commandline.WriteString(opts.RecipeURL)

	for key, value := range values {
		commandline.WriteString(fmt.Sprintf(" --set \"%s=%s\"", key, value))
	}
	retryMessage := fmt.Sprintf("To re-run the recipe with the same values, use the following command:\n\n%s", commandline.String())
	return retryMessage
}
