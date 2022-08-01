package main

import (
	"fmt"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var validateCmd = &cobra.Command{
		Use:          "validate",
		Short:        "validate a recipe",
		RunE:         validateFunc,
		SilenceUsage: true,
	}

	return validateCmd
}

func validateFunc(cmd *cobra.Command, args []string) error {
	r, err := recipe.Load(args[0])
	if err != nil {
		return fmt.Errorf("Error when loading the recipe: %s\n", err)
	}

	err = r.Validate()
	if err != nil {
		return fmt.Errorf("Validation failed: %s\n", err)
	}

	return nil
}
