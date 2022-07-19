package main

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
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

	err = r.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	values, _ := promptUserForValues(r.Variables)

	context := map[string]interface{}{
		"Variables": values,
	}

	output, err := engine.Render(r, context)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("OUTPUT: %v", output)

	// TODO: Write output files to the given path
}

func promptUserForValues(variables recipe.VariableMap) (recipe.VariableValues, error) {
	values := recipe.VariableValues{}
	for name, variable := range variables { // TODO: Sort variables before looping for consistent behaviour
		var prompt survey.Prompt

		if len(variable.Options) != 0 {
			prompt = &survey.Select{
				Message: name,
				Default: variable.Default,
				Help:    variable.Description,
				Options: variable.Options,
			}
		} else {
			prompt = &survey.Input{
				Message: name,
				Default: variable.Default,
				Help:    variable.Description,
			}
		}

		opts := make([]survey.AskOpt, 0)

		if !variable.Optional {
			opts = append(opts, survey.WithValidator(survey.Required))
		}

		if variable.RegExp.Pattern != "" {
			validator, err := variable.RegExp.CreateValidatorFunc()
			if err != nil {
				return recipe.VariableValues{}, err
			}

			opts = append(opts, survey.WithValidator(validator))
		}

		var answer string
		err := survey.AskOne(prompt, &answer, opts...)
		if err != nil {
			return recipe.VariableValues{}, err
		}

		values[name] = answer
	}

	return values, nil
}
