package main

import (
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newExecuteCmd() *cobra.Command {
	// execCmd represents the exec command
	var execCmd = &cobra.Command{
		Use:     "execute",
		Aliases: []string{"exec"},
		Short:   "Execute a given recipe",
		Run:     executeFunc,
	}

	return execCmd
}

func executeFunc(cmd *cobra.Command, args []string) {
	// TODO: Read recipe from directory

	// TODO: Prompt user to fill the variables

	values := map[string]interface{}{
		"Variables": recipe.VariableValues{
			"PROJECT_NAME": "my-project",
		},
	}
	output, _ := engine.Render(&recipe.Recipe{}, values)

	fmt.Println(output)
}
