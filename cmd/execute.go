package main

import (
	"errors"
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

func newExecuteCmd() *cobra.Command {
	// execCmd represents the exec command
	var execCmd = &cobra.Command{
		Use:     "execute",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a path argument")
			}
			return nil
		},
		Run: executeFunc,
	}

	return execCmd
}

func executeFunc(cmd *cobra.Command, args []string) {
	r, err := recipe.Load(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Prompt user to fill the variables

	values := map[string]interface{}{
		"Variables": recipe.VariableValues{
			"PROJECT_NAME": "my-project",
		},
	}
	output, err := engine.Render(r, values)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("OUTPUT: %v", output)
}
