package recipeutil

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
)

const ValueEnvVarPrefix = "JALAPENO_VAR_"

var (
	ErrVarNotDefinedInRecipe = errors.New("following variable does not exist in the recipe")
)

func ParsePredefinedValues(variables []recipe.Variable, flags []string) (recipe.VariableValues, error) {
	values := make(recipe.VariableValues)
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, ValueEnvVarPrefix) {
			continue
		}

		// Add environment variables at the beginning of the slice so CLI flags override env. variables
		flags = append([]string{strings.TrimPrefix(env, ValueEnvVarPrefix)}, flags...)
	}

	for _, flag := range flags {
		splitted := strings.SplitN(flag, "=", 2)
		if len(splitted) != 2 {
			return nil, fmt.Errorf("TODO %s", flag)
		}
		varName := splitted[0]
		varValue := splitted[1]

		if varValue == "true" {
			values[varName] = true
		} else if varValue == "false" {
			values[varName] = false
		} else {
			values[varName] = varValue
		}
	}

	for varName, value := range values {
		found := false
		for _, variable := range variables {
			if variable.Name != varName {
				continue
			}

			found = true
			if variable.RegExp.Pattern != "" {
				validator := variable.RegExp.CreateValidatorFunc()
				if err := validator(value); err != nil {
					return nil, fmt.Errorf("validator failed for value '%s=%s': %w", varName, value, err)
				}
			}
			break
		}

		if !found {
			return nil, fmt.Errorf("%w: %s", ErrVarNotDefinedInRecipe, varName)
		}
	}

	return values, nil
}

func MergeValues(valuesSlice ...recipe.VariableValues) recipe.VariableValues {
	merged := make(recipe.VariableValues)
	for _, values := range valuesSlice {
		for key := range values {
			merged[key] = values[key]
		}
	}

	return merged
}

func FilterVariables(vars []recipe.Variable, values recipe.VariableValues) []recipe.Variable {
	variablesWithoutValues := make([]recipe.Variable, 0, len(vars))
	for _, variable := range vars {
		if _, exists := values[variable.Name]; !exists {
			variablesWithoutValues = append(variablesWithoutValues, variable)
		}
	}

	return variablesWithoutValues
}
