package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

var (
	outputBasePath = ""
)

func newExecuteCmd() *cobra.Command {
	var execCmd = &cobra.Command{
		Use:     "execute",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe and save output to a path",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a path argument")
			}
			return nil
		},
		Run: executeFunc,
	}

	execCmd.Flags().StringVarP(&outputBasePath, "output", "o", ".", "Path where the output files should be created")

	return execCmd
}

func executeFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(outputBasePath); os.IsNotExist(err) {
		fmt.Println("Output path does not exist")
		return
	}

	re, err := recipe.Load(args[0])
	if err != nil {
		fmt.Printf("Error when loading the recipe: %s\n", err)
		return
	}

	fmt.Printf("Recipe name: %s\n\n", re.Metadata.Name)

	err = re.Validate()
	if err != nil {
		fmt.Printf("The provided recipe was invalid: %s\n", err)
		return
	}

	values, err := promptUserForValues(re.Variables)
	if err != nil {
		fmt.Printf("Error when prompting for values: %s\n", err)
		return
	}

	// Define the context which is available on templates
	context := map[string]interface{}{
		"Variables": values,
	}

	renderedFiles, err := engine.Render(re, context)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create sub directory for recipe
	recipePath := filepath.Join(outputBasePath, ".jalapeno")
	err = os.MkdirAll(recipePath, 0700)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = recipeutil.SaveRecipe(re, recipePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = recipeutil.SaveFiles(renderedFiles, outputBasePath)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func promptUserForValues(variables []recipe.Variable) (recipe.VariableValues, error) {
	values := recipe.VariableValues{}

	for _, variable := range variables {
		// TODO: Check if the value for the variable was alredy provided by CLI arguments

		var prompt survey.Prompt

		if len(variable.Options) != 0 {
			prompt = &survey.Select{
				Message: variable.Name,
				Help:    variable.Description,
				Options: variable.Options,
			}
		} else {
			prompt = &survey.Input{
				Message: variable.Name,
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

		if _, exist := values[variable.Name]; exist {
			return recipe.VariableValues{}, fmt.Errorf(`variable "%s" has been declared multiple times`, variable.Name)
		}

		values[variable.Name] = answer
	}

	return values, nil
}
