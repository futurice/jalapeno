package main

import (
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create NAME",
		Short: "Create a new recipe",
		Long: `
...
	foo/
	├── recipe.yml
	├── templates/
`, // TODO
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
		cmd.PrintErrf("directory '%s' already exists", recipeName)
		return
	}

	err := os.Mkdir(path, 0700)
	if err != nil {
		cmd.PrintErrf("can not create directory %s: %v", path, err)
		return
	}

	err = re.Validate()
	if err != nil {
		cmd.PrintErrln("internal error: placeholder recipe is not valid")
		return
	}

	err = re.Save(path)
	if err != nil {
		cmd.PrintErrf("can not save recipe to the directory: %v", err)
		return
	}

	err = os.Mkdir(filepath.Join(path, recipe.RecipeTemplatesDirName), 0700)
	if err != nil {
		cmd.PrintErrf("can not save templates to the directory: %v", err)
		return
	}
}
