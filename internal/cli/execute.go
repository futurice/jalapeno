package cli

import (
	"os"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

type executeOptions struct {
	RecipePath string
	option.Values
	option.Output
	option.Common
}

func NewExecuteCmd() *cobra.Command {
	var opts executeOptions
	var cmd = &cobra.Command{
		Use:     "execute RECIPE_PATH",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe and save output to path",
		Long:    "TODO",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runExecute(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runExecute(cmd *cobra.Command, opts executeOptions) {
	if _, err := os.Stat(opts.OutputPath); os.IsNotExist(err) {
		cmd.PrintErrln("Error: output path does not exist")
		return
	}

	re, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	cmd.Printf("Recipe name: %s\n", re.Metadata.Name)

	if re.Metadata.Description != "" {
		cmd.Printf("Description: %s\n", re.Metadata.Description)
	}

	predefinedValues, err := recipeutil.ParsePredefinedValues(re.Variables, opts.Values.Flags)
	if err != nil {
		cmd.PrintErrf("Error when parsing provided values: %v\n", err)
		return
	}

	values, err := recipeutil.PromptUserForValues(recipeutil.FilterVariables(re.Variables, predefinedValues))
	if err != nil {
		cmd.PrintErrf("Error when prompting for values: %v\n", err)
		return
	}

	sauce, err := re.Execute(
		recipeutil.MergeValues(values, predefinedValues),
		uuid.Must(uuid.NewV4()),
	)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	// Load all existing sauces
	existingSauces, err := recipe.LoadSauces(opts.OutputPath)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	// Check for conflicts
	for _, s := range existingSauces {
		if conflicts := s.Conflicts(sauce); conflicts != nil {
			cmd.PrintErrf("conflict in recipe '%s': file '%s' was already created by recipe '%s'\n", re.Name, conflicts[0].Path, s.Recipe.Name)
			return
		}
	}

	err = sauce.Save(opts.OutputPath)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	cmd.Println("\nRecipe executed successfully!")

	if re.InitHelp != "" {
		cmd.Printf("Next up: %s\n", re.InitHelp)
	}
}
