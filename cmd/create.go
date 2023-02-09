package main

import (
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type createOptions struct {
	RecipeName string
	option.Common
}

func newCreateCmd() *cobra.Command {
	var opts createOptions
	// createCmd represents the create command
	var cmd = &cobra.Command{
		Use:   "create NAME",
		Short: "Create a new recipe",
		Long: `
...
	foo/
	├── recipe.yml
	├── templates/
`, // TODO
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeName = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCreate(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCreate(cmd *cobra.Command, opts createOptions) {
	re := createExampleRecipe(opts.RecipeName)

	path := filepath.Join(".", opts.RecipeName)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		cmd.PrintErrf("directory '%s' already exists\n", opts.RecipeName)
		return
	}

	err := os.Mkdir(path, 0700)
	if err != nil {
		cmd.PrintErrf("can not create directory %s: %v\n", path, err)
		return
	}

	err = re.Validate()
	if err != nil {
		// TODO: Clean up the already create directory if the command failed
		cmd.PrintErrln("internal error: placeholder recipe is not valid")
		return
	}

	err = re.Save(path)
	if err != nil {
		// TODO: Clean up the already create directory if the command failed
		cmd.PrintErrf("can not save recipe to the directory: %v\n", err)
		return
	}

	err = os.Mkdir(filepath.Join(path, recipe.RecipeTemplatesDirName), 0700)
	if err != nil {
		// TODO: Clean up the already create directory if the command failed
		cmd.PrintErrf("can not save templates to the directory: %v\n", err)
		return
	}
}

func createExampleRecipe(name string) *recipe.Recipe {
	return &recipe.Recipe{
		Metadata: recipe.Metadata{
			APIVersion:  "v1",
			Name:        name,
			Version:     "v0.0.0",
			Description: "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools",
		},
		Variables: []recipe.Variable{
			{Name: "MY_VAR", Default: "Hello World!"},
		},
	}
}
