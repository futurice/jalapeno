package main

import (
	"os"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new recipe",
		Run:   createFunc,
	}

	createCmd.Flags().StringVarP(&outputBasePath, "output", "o", ".", "Path where the output files should be created")

	return createCmd
}

func createFunc(cmd *cobra.Command, args []string) {
	name := "example" // TODO: Get from arguments
	r := &recipe.Recipe{
		Metadata: recipe.Metadata{
			Name:        name,
			Version:     "0.0.0",
			Description: "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools",
		},
	}

	err := r.Validate()
	if err != nil {
		panic("invalid example recipe") // TODO
	}

	recipeutil.SaveRecipe(r, ".")

	err = os.Mkdir("./templates", 0700)
	if err != nil {
		panic("can not create directory") // TODO
	}
}
