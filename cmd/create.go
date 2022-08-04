package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create <recipe_name>",
		Short: "Create a new recipe",
		Long: `
		...
			foo/
			├── recipe.yml
			├── templates/
		`,
		Args: cobra.ExactArgs(1),
		Run:  createFunc,
	}

	return createCmd
}

func createFunc(cmd *cobra.Command, args []string) {
	recipeName := args[0]
	re := &recipe.Recipe{
		Metadata: recipe.Metadata{
			Name:        recipeName,
			Version:     "0.0.0",
			Description: "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools",
		},
		Variables: []recipe.Variable{
			{Name: "MY_VAR", Default: "Hello World!"},
		},
	}

	path := filepath.Join(".", recipeName)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Printf("error: directory '%s' already exists", recipeName)
		return
	}

	err := os.Mkdir(path, 0700)
	if err != nil {
		panic("can not create directory") // TODO
	}

	err = re.Validate()
	if err != nil {
		panic("invalid example recipe") // TODO
	}

	err = re.Save(path)
	if err != nil {
		panic("recipe saving failed") // TODO
	}

	err = os.Mkdir(filepath.Join(path, recipe.RecipeTemplatesDirName), 0700)
	if err != nil {
		panic("can not create directory") // TODO
	}
}
