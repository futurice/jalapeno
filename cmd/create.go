package main

import (
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new recipe",
		Run:   createFunc,
	}

	return createCmd
}

func createFunc(cmd *cobra.Command, args []string) {
	name := "example" // TODO: Get from arguments
	metadata := recipe.Metadata{
		Name:    name,
		Version: "0.0.0",
	}

	metadata.Validate()
}
