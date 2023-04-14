package main

import (
	"os"

	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

type executeOptions struct {
	RecipePath string
	option.Output
	option.Common
}

func newExecuteCmd() *cobra.Command {
	var opts executeOptions
	var cmd = &cobra.Command{
		Use:     "execute RECIPE",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe and save output to path",
		Long:    "", // TODO
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

	// TODO: Set values provided by --set flag to re.Values

	values, err := recipeutil.PromptUserForValues(re.Variables)
	if err != nil {
		cmd.PrintErrf("Error when prompting for values: %v\n", err)
		return
	}

	sauce, err := re.Execute(values, uuid.Must(uuid.NewV4()))
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
		conflicts := s.Conflicts(sauce)
		if conflicts != nil {
			cmd.PrintErrf("conflict in recipe %s: %s was already created by recipe %s\n", re.Name, conflicts[0].Path, s.Recipe.Name)
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
