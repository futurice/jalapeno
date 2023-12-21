package recipeutil

import (
	"encoding/csv"
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

func ParseProvidedValues(variables []recipe.Variable, flags []string, delimiter rune) (recipe.VariableValues, error) {
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
			return nil, fmt.Errorf("invalid format on flag '%s'. Use format 'MY_VAR=VALUE'", flag)
		}
		varName := splitted[0]
		varValue := splitted[1]

		var targetedVariable *recipe.Variable
		for _, variable := range variables {
			if variable.Name != varName {
				continue
			} else {
				targetedVariable = &variable
				break
			}
		}

		if targetedVariable == nil {
			return nil, fmt.Errorf("%w: %s", ErrVarNotDefinedInRecipe, varName)
		}

		for i := range targetedVariable.Validators {
			validatorFunc := targetedVariable.Validators[i].CreateValidatorFunc()
			if err := validatorFunc(varValue); err != nil {
				return nil, fmt.Errorf("validator failed for value '%s=%s': %w", varName, varValue, err)
			}
		}

		switch {
		case targetedVariable.Confirm:
			if varValue == "true" {
				values[varName] = true
			} else if varValue == "false" {
				values[varName] = false
			} else {
				return nil, fmt.Errorf("value provided for variable '%s' was not a boolean", varName)
			}
		case len(targetedVariable.Columns) > 0:
			varValue = strings.ReplaceAll(varValue, "\\n", "\n")
			table, err := CSVToTable(targetedVariable.Columns, varValue, delimiter)
			if err != nil {
				return nil, fmt.Errorf("failed to parse table from CSV for variable '%s': %w", varName, err)
			}
			values[varName] = table

		default:
			values[varName] = varValue
		}
	}

	return values, nil
}

// MergeValues merges multiple VariableValues into one. If a key exists in multiple VariableValues, the value from the
// last VariableValues will be used.
func MergeValues(valuesSlice ...recipe.VariableValues) recipe.VariableValues {
	merged := make(recipe.VariableValues)
	for _, values := range valuesSlice {
		for key := range values {
			merged[key] = values[key]
		}
	}

	return merged
}

func FilterVariablesWithoutValues(variables []recipe.Variable, values recipe.VariableValues) []recipe.Variable {
	variablesWithoutValues := make([]recipe.Variable, 0, len(variables))
	for _, variable := range variables {
		if _, exists := values[variable.Name]; !exists {
			variablesWithoutValues = append(variablesWithoutValues, variable)
		}
	}

	return variablesWithoutValues
}

func CSVToTable(columns []string, str string, delimiter rune) ([]map[string]string, error) {
	reader := csv.NewReader(strings.NewReader(str))
	reader.FieldsPerRecord = len(columns)
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	table := make([]map[string]string, len(rows))
	for i, row := range rows {
		table[i] = make(map[string]string)
		for j, cell := range row {
			table[i][columns[j]] = cell
		}
	}

	return table, nil
}

func TableToCSV(table []map[string]string, delimiter rune) (string, error) {
	// declare an io writer that writes to a string
	var stringWriter strings.Builder
	csvWriter := csv.NewWriter(&stringWriter)
	csvWriter.Comma = delimiter
	for _, row := range table {
		var csvRow []string
		for _, column := range row {
			csvRow = append(csvRow, column)
		}
		if err := csvWriter.Write(csvRow); err != nil {
			return "", err
		}
	}
	csvWriter.Flush()
	return stringWriter.String(), nil
}

func RowsToTable(columns []string, rows [][]string) ([]map[string]string, error) {
	table := make([]map[string]string, len(rows))
	for i, row := range rows {
		table[i] = make(map[string]string)
		for j, cell := range row {
			table[i][columns[j]] = cell
		}
	}

	return table, nil
}

func NewNoInputError(vars []recipe.Variable) error {
	var errMsg string
	switch len(vars) {
	case 0:
		return errors.New("there was file conflicts which needs to be manually resolved while `--no-input` flag was set to true")
	case 1:
		return fmt.Errorf("value for variable %s is missing and `--no-input` flag was set to true", vars[0].Name)
	default:
		varNames := make([]string, len(vars))
		for i, v := range vars {
			varNames[i] = v.Name
		}
		errMsg = fmt.Sprintf("values for variables [%s] are", strings.Join(varNames, ","))
		return fmt.Errorf("%s missing and `--no-input` flag was set to true", errMsg)
	}
}
