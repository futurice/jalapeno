package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
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
	r, err := recipe.Load(args[0])
	if err != nil {
		fmt.Printf("Error when loading the recipe: %s\n", err)
		return
	}

	err = r.Validate()
	if err != nil {
		fmt.Printf("The provided recipe was invalid: %s\n", err)
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

	err = writeMapToFiles(output, outputBasePath)
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

func writeMapToFiles(files map[string]string, basepath string) error {
	if _, err := os.Stat(basepath); os.IsNotExist(err) {
		return err
	}

	for filename, data := range files {
		path := filepath.Join(basepath, filename)

		// Create file's parent directories (if not already exist)
		err := os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}

		// Create the file
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write the data to the file
		_, err = f.WriteString(data)
		if err != nil {
			return err
		}

		f.Sync()
	}
	return nil
}
