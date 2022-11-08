package main

import (
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var validateCmd = &cobra.Command{
		Use:   "validate RECIPE",
		Short: "Validate a recipe",
		Long:  "", // TODO
		Run:   validateFunc,
		Args:  cobra.ExactArgs(1),
	}

	return validateCmd
}

func validateFunc(cmd *cobra.Command, args []string) {
	r, err := recipe.Load(args[0])
	if err != nil {
		cmd.PrintErrf("could not load the recipe: %s\n", err)
		return
	}

	err = r.Validate()
	if err != nil {
		cmd.PrintErrf("validation failed: %s\n", err)
		return
	}

	cmd.Println("Validation ok")
}
