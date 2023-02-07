package main

import (
	"os"

	"github.com/futurice/jalapeno/internal/option"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
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

	option.ApplyFlags(&opts, cmd.Flags())

	return cmd
}

func runExecute(cmd *cobra.Command, opts executeOptions) {
	if _, err := os.Stat(opts.OutputPath); os.IsNotExist(err) {
		cmd.PrintErrln("output path does not exist")
		return
	}

	re, err := recipe.Load(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("can't load the recipe: %v\n", err)
		return
	}

	cmd.Printf("Recipe name: %s\n", re.Metadata.Name)

	if re.Metadata.Description != "" {
		cmd.Printf("Description: %s\n", re.Metadata.Description)
	}

	err = re.Validate()
	if err != nil {
		cmd.PrintErrf("the provided recipe was invalid: %v\n", err)
		return
	}

	if len(re.Templates) == 0 {
		cmd.PrintErrf("the recipe does not contain any templates\n")
		return
	}

	// TODO: Set values provided by --set flag to re.Values

	err = recipeutil.PromptUserForValues(re)
	if err != nil {
		cmd.PrintErrf("error when prompting for values: %v\n", err)
		return
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// Load all rendered recipes
	rendered, err := recipe.LoadRendered(opts.OutputPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// Check for conflicts
	for _, r := range rendered {
		conflicts := re.Conflicts(&r)
		if conflicts != nil {
			cmd.PrintErrf("conflict in recipe %s: %s was already created by recipe %s\n", re.Name, conflicts[0].Path, r.Name)
			return
		}
	}

	err = re.Save(opts.OutputPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = recipeutil.SaveFiles(re.Files, opts.OutputPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("\nRecipe executed successfully!")

	if re.InitHelp != "" {
		cmd.Printf("Next up: %s\n", re.InitHelp)
	}
}
