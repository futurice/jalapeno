package main

import (
	"github.com/futurice/jalapeno/cmd/internal/option"
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
	├──── README.md
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

	err := re.Validate()
	if err != nil {
		cmd.PrintErrln("Internal error: placeholder recipe is not valid")
		return
	}

	err = re.Save(".")
	if err != nil {
		cmd.PrintErrf("Error: can not save recipe to the directory: %v\n", err)
		return
	}
}

func createExampleRecipe(name string) *recipe.Recipe {
	r := recipe.NewRecipe()
	r.Metadata.Name = name
	r.Metadata.Version = "v0.0.0"
	r.Metadata.Description = "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools"
	r.Variables = []recipe.Variable{
		{Name: "MY_VAR", Default: "Hello World!"},
	}
	r.Templates = map[string][]byte{
		"README.md": []byte("{{ .Variables.MY_VAR }}"),
	}

	return r
}
