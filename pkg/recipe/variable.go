package recipe

import (
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/expr-lang/expr"
	"gopkg.in/yaml.v3"
)

type Variable struct {
	// The name of the variable. It is also used as unique identifier.
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`

	// Default value of the variable
	Default string `yaml:"default,omitempty"`

	// If set to true, the prompt will be yes/no question, and the value type will be boolean
	Confirm bool `yaml:"confirm,omitempty"`

	// If set to true, the variable can be left empty
	Optional bool `yaml:"optional,omitempty"`

	// The user selects the value from a list of options
	Options []string `yaml:"options,omitempty"`

	// If set to true, the user can select multiple values from the list of options
	Multi bool `yaml:"multi,omitempty"`

	// Validators for the variable
	Validators []VariableValidator `yaml:"validators,omitempty"`

	// Makes the variable conditional based on the result of the expression. The result of the evaluation needs to be a boolean value. Uses https://github.com/expr-lang/expr
	If string `yaml:"if,omitempty"`

	// Set the variable as a table type with columns defined by this property
	Columns []string `yaml:"columns,omitempty"`
}

type VariableType uint8

const (
	VariableTypeUnknown VariableType = iota
	VariableTypeString
	VariableTypeBoolean
	VariableTypeSelect
	VariableTypeMultiSelect
	VariableTypeTable
)

type VariableValidator struct {
	// Regular expression pattern to match the input against
	Pattern string `yaml:"pattern,omitempty"`

	// If the regular expression validation fails, this help message will be shown to the user
	Help string `yaml:"help,omitempty"`

	// Apply the validator to a column if the variable type is table
	Column string `yaml:"column,omitempty"`

	// When targeting table columns, set this to true to make sure that the values in the column are unique
	Unique bool `yaml:"unique,omitempty"`
}

// VariableValues stores values for each variable
type VariableValues map[string]any

type MultiSelectValue []string

type TableValue struct {
	Columns []string   `yaml:"columns"`
	Rows    [][]string `yaml:"rows,flow"`
}

var _ yaml.Unmarshaler = (*VariableValues)(nil)

var startsWithNumber = regexp.MustCompile(`^\d.*`)

func (v *Variable) Validate() error {
	if v.Name == "" {
		return errors.New("variable name is required")
	}

	if startsWithNumber.MatchString(v.Name) {
		return errors.New("variable name can not start with a number")
	}

	if v.DetermineType() == VariableTypeUnknown {
		return errors.New("internal error: variable type could not be determined")
	}

	if v.Confirm {
		if len(v.Options) > 0 {
			return errors.New("`confirm` and `options` properties can not be defined at the same time")
		} else if len(v.Columns) > 0 {
			return errors.New("`confirm` and `columns` properties can not be defined at the same time")
		} else if v.Optional {
			return errors.New("boolean variables can not be optional")
		}
	}

	if len(v.Options) > 0 {
		if len(v.Columns) > 0 {
			return errors.New("`options` and `columns` properties can not be defined at the same time")
		} else if v.Optional && !v.Multi {
			return errors.New("select variables can not be optional")
		}
	}

	if v.Multi && len(v.Options) == 0 {
		return errors.New("multiselect variables need to have options defined")
	}

	for i, validator := range v.Validators {
		validatorIndex := fmt.Sprintf("validator %d", i+1)
		if v.Confirm {
			return fmt.Errorf("%s: validators for boolean variables are not supported", validatorIndex)
		}

		if len(v.Options) > 0 {
			return fmt.Errorf("%s: validators for select variables are not supported", validatorIndex)
		}

		if len(v.Columns) > 0 && validator.Column == "" {
			return fmt.Errorf("%s: validator need to have `column` property defined since the variable is table type", validatorIndex)
		}

		if validator.Unique {
			if validator.Column == "" {
				return fmt.Errorf("%s: validator need to have `column` property defined since unique validation works only on table variables", validatorIndex)
			}
			if validator.Pattern != "" {
				return fmt.Errorf("%s: validator can not have `pattern` property defined when `unique` is set to true", validatorIndex)
			}
			return nil
		} else {
			if validator.Pattern == "" {
				return fmt.Errorf("%s: regexp pattern is empty", validatorIndex)
			}
			if _, err := regexp.Compile(validator.Pattern); err != nil {
				return fmt.Errorf("%s: invalid validator regexp pattern: %w", validatorIndex, err)
			}
		}

		if validator.Column != "" {
			if len(v.Columns) == 0 {
				return fmt.Errorf("%s: validator is defined for column while the variable has not defined any", validatorIndex)
			}

			found := false
			for _, c := range v.Columns {
				if c == validator.Column {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("%s: column %s does not exist in the variable", validatorIndex, validator.Column)
			}
		}
	}

	if v.If != "" {
		if _, err := expr.Compile(v.If); err != nil {
			return fmt.Errorf("invalid 'if' expression: %w", err)
		}
	}

	return nil
}

// NOTE: This function does note validate the values against the variable definitions.
// It only checks if the name of the values are not empty and the values are of supported types.
func (val VariableValues) Validate() error {
	for name, v := range val {
		if name == "" {
			return errors.New("variable name can not be empty")
		}

		switch v.(type) {
		case string, bool, MultiSelectValue, TableValue:
			break
		default:
			return fmt.Errorf("unsupported variable value type (%T)", v)
		}
	}

	return nil
}

// RequiresTableContext returns true if the validator function should be created with CreateTableValidatorFunc
func (r *VariableValidator) RequiresTableContext() bool {
	return r.Unique
}

func (r *VariableValidator) CreateTableValidatorFunc() (func(cols []string, rows [][]string, input string) error, error) {
	if r.Unique {
		return func(cols []string, rows [][]string, input string) error {
			colIndex := slices.Index(cols, r.Column)
			colValues := make([]string, len(rows))
			for i, row := range rows {
				colValues[i] = row[colIndex]
			}
			slices.Sort(colValues)

			if uniqValues := len(slices.Compact(colValues)); uniqValues != len(colValues) {
				if r.Help != "" {
					return errors.New(r.Help)
				} else {
					return errors.New("value not unique within column")
				}
			}

			return nil
		}, nil
	}

	return nil, fmt.Errorf("unsupported table validator on column %q", r.Column)
}

func (r *VariableValidator) CreateValidatorFunc() (func(input string) error, error) {
	if r.Pattern != "" {
		reg := regexp.MustCompile(r.Pattern)

		return func(input string) error {
			if match := reg.MatchString(input); !match {
				if r.Help != "" {
					return errors.New(r.Help)
				} else {
					return errors.New("the input did not match the regexp pattern")
				}
			}
			return nil
		}, nil
	}

	return nil, fmt.Errorf("unsupported validator on column %q", r.Column)
}

func (t *TableValue) FromCSV(columns []string, input string, delimiter rune) error {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = len(columns)
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	t.Columns = columns
	t.Rows = make([][]string, len(rows))
	for i, row := range rows {
		t.Rows[i] = make([]string, len(columns))
		copy(t.Rows[i], row)
	}

	return nil
}

func (t TableValue) ToCSV(delimiter rune) (string, error) {
	var stringWriter strings.Builder
	csvWriter := csv.NewWriter(&stringWriter)
	csvWriter.Comma = delimiter

	for _, row := range t.Rows {
		csvRow := make([]string, len(t.Columns))
		for i := range t.Columns {
			csvRow[i] = row[i]
		}

		if err := csvWriter.Write(csvRow); err != nil {
			return "", err
		}
	}
	csvWriter.Flush()
	return stringWriter.String(), nil
}

func (t TableValue) ToMapSlice() []map[string]string {
	output := make([]map[string]string, len(t.Rows))
	for i, row := range t.Rows {
		output[i] = make(map[string]string, len(t.Columns))
		for j, column := range t.Columns {
			output[i][column] = row[j]
		}
	}

	return output
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (vv *VariableValues) UnmarshalYAML(value *yaml.Node) error {
	rawYaml := make(map[string]any)
	err := value.Decode(rawYaml)
	if err != nil {
		return err
	}

	*vv = make(VariableValues, len(rawYaml))
	for name, value := range rawYaml {
		switch v := value.(type) {

		// Check if the value is a TableValue
		case map[string]any:
			_, columnsExist := v["columns"]
			_, rowsExist := v["rows"]

			// If the value is a TableValue, parse it
			if columnsExist && rowsExist {
				rawColumns, ok := v["columns"].([]any)
				if !ok {
					return fmt.Errorf("failed to parse table columns for variable '%s'", name)
				}

				columns := make([]string, len(rawColumns))
				for i, c := range rawColumns {
					columns[i], ok = c.(string)
					if !ok {
						return fmt.Errorf("failed to parse table column for variable '%s'", name)
					}
				}

				rawRows, ok := v["rows"].([]any)
				if !ok {
					return fmt.Errorf("failed to parse table rows for variable '%s'", name)
				}

				rows := make([][]string, len(rawRows))
				for i := range rawRows {
					rawRow, ok := rawRows[i].([]any)
					if !ok {
						return fmt.Errorf("failed to parse table row for variable '%s'", name)
					}

					rows[i] = make([]string, len(rawRow))
					for j, c := range rawRow {
						rows[i][j], ok = c.(string)
						if !ok {
							return fmt.Errorf("failed to parse table cell for variable '%s'", name)
						}
					}
				}

				(*vv)[name] = TableValue{
					Columns: columns,
					Rows:    rows,
				}
			}
		// Multiselect values without a value are an empty []any and we want to cast them into MultiSelectValues
		case []any:
			multiV := make(MultiSelectValue, len(v))
			for i, interfaceV := range v {
				str, ok := interfaceV.(string)
				if !ok {
					return fmt.Errorf("failed to read value to multi select value: %v", v)
				}
				multiV[i] = str
			}
			(*vv)[name] = multiV
		default:
			(*vv)[name] = v
		}
	}

	return nil
}

func (v Variable) DetermineType() VariableType {
	switch {
	case v.Confirm:
		return VariableTypeBoolean
	case len(v.Options) > 0:
		if v.Multi {
			return VariableTypeMultiSelect
		} else {
			return VariableTypeSelect
		}
	case len(v.Columns) > 0:
		return VariableTypeTable
	default:
		return VariableTypeString
	}
}

func (v Variable) ParseDefaultValue() (any, error) {
	switch v.DetermineType() {
	case VariableTypeBoolean:
		return v.Default == "true", nil
	case VariableTypeSelect:
		return v.Default, nil
	case VariableTypeMultiSelect:
		return strings.Split(v.Default, ","), nil
	case VariableTypeTable:
		t := TableValue{}
		err := t.FromCSV(v.Columns, v.Default, ',')
		if err != nil {
			return nil, err
		}
		return t, nil
	case VariableTypeString:
		return v.Default, nil
	default:
		return nil, errors.New("unknown variable type")
	}
}

func (v MultiSelectValue) ToString(delimiter rune) string {
	return strings.Join(v, string(delimiter))
}

func (v *MultiSelectValue) FromString(s string, delimiter rune) {
	*v = strings.Split(s, string(delimiter))
}
