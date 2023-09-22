package recipeutil

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/antonmedv/expr"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func PromptUserForValues(variables []recipe.Variable, existingValues recipe.VariableValues) (recipe.VariableValues, error) {
	// TODO: This command does not respect stdio defined by the Cobra cmd, so
	// capturing and examining the output of this function does not work at the moment
	values := recipe.VariableValues{}
	headerAdded := false

	for _, variable := range variables {
		if !headerAdded {
			fmt.Println("\nProvide the following variables:")
			headerAdded = true
		}

		var prompt survey.Prompt
		var askFunc AskFunc = askString

		if variable.If != "" {
			result, err := expr.Eval(variable.If, MergeValues(existingValues, values))
			if err != nil {
				return nil, fmt.Errorf("error when evaluating 'if' expression: %w", err)
			}

			variableShouldBePrompted, ok := result.(bool)
			if !ok {
				return nil, fmt.Errorf("result of 'if' expression was not a boolean value, was %T instead", result)
			}

			if !variableShouldBePrompted {
				continue
			}
		}

		// Select with predefined options
		if len(variable.Options) != 0 {
			prompt = &survey.Select{
				Message: variable.Name,
				Help:    variable.Description,
				Options: variable.Options,
			}

			// Yes/No question
		} else if variable.Confirm {
			prompt = &survey.Confirm{
				Message: variable.Name,
				Help:    variable.Description,
				Default: variable.Default == "true",
			}
			askFunc = askBool

			// Free input question
		} else {
			prompt = &survey.Input{
				Message: variable.Name,
				Default: variable.Default,
				Help:    variable.Description,
			}
		}

		opts := make([]survey.AskOpt, 0)

		if !(variable.Optional || variable.If != "") {
			opts = append(opts, survey.WithValidator(survey.Required))
		}

		if variable.RegExp.Pattern != "" {
			validator := variable.RegExp.CreateValidatorFunc()
			opts = append(opts, survey.WithValidator(validator))
		}

		answer, err := askFunc(prompt, opts)
		if err != nil {
			return nil, err
		}

		values[variable.Name] = answer
	}

	return values, nil
}

// NOTE: Since survey.AskOne tries to cast the answer to the type of the response
// value pointer and the type of response value can not be interface{},
// we need to create different ask functions for each response type and return interface{}
type AskFunc func(prompt survey.Prompt, opts []survey.AskOpt) (interface{}, error)

func askString(prompt survey.Prompt, opts []survey.AskOpt) (interface{}, error) {
	return ask[string](prompt, opts)
}

func askBool(prompt survey.Prompt, opts []survey.AskOpt) (interface{}, error) {
	return ask[bool](prompt, opts)
}

func ask[T string | bool](prompt survey.Prompt, opts []survey.AskOpt) (T, error) {
	var answer T
	if err := survey.AskOne(prompt, &answer, opts...); err != nil {
		return answer, err
	}
	return answer, nil
}
