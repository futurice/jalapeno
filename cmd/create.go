package main

import (
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

// Any type can be given to the select's item as long as the templates properly implement the dot notation
// to display it.
type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

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
	metadata := recipe.Metadata{}

	metadata.Validate()
}
