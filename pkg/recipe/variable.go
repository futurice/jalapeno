package recipe

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/antonmedv/expr"
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

	// Validators for the variable
	Validators []VariableValidator `yaml:"validators,omitempty"`

	// Makes the variable conditional based on the result of the expression. The result of the evaluation needs to be a boolean value. Uses https://github.com/antonmedv/expr
	If string `yaml:"if,omitempty"`

	// Set the variable as a table type with columns defined by this property
	Columns []string `yaml:"columns,omitempty"`
}

type VariableValidator struct {
	// Regular expression pattern to match the input against
	Pattern string `yaml:"pattern,omitempty"`

	// If the regular expression validation fails, this help message will be shown to the user
	Help string `yaml:"help,omitempty"`

	// Apply the validator to a column if the variable type is table
	Column string `yaml:"column,omitempty"`
}

// VariableValues stores values for each variable
type VariableValues map[string]interface{}

func (v *Variable) Validate() error {
	if v.Name == "" {
		return errors.New("variable name is required")
	}

	if v.Confirm {
		if len(v.Options) > 0 {
			return errors.New("`confirm` and `options` properties can not be defined at the same time")
		} else if len(v.Columns) > 0 {
			return errors.New("`confirm` and `columns` properties can not be defined at the same time")
		}
	}

	for i, validator := range v.Validators {
		baseErr := fmt.Errorf("validator %d", i+1)
		if v.Confirm {
			return fmt.Errorf("%s: validators for boolean variables are not supported", baseErr)
		}

		if len(v.Options) > 0 {
			return fmt.Errorf("%s: validators for select variables are not supported", baseErr)
		}

		if len(v.Columns) > 0 && validator.Column == "" {
			return fmt.Errorf("%s: validator need to have `column` property defined since the variable is table type", baseErr)
		}

		if validator.Pattern == "" {
			return fmt.Errorf("%s: regexp pattern is empty", baseErr)
		}

		if validator.Column != "" {
			if len(v.Columns) == 0 {
				return fmt.Errorf("%s: validator is defined for column while the variable has not defined any", baseErr)
			}

			found := false
			for _, c := range v.Columns {
				if c == validator.Column {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("%s: column %s does not exist in the variable", baseErr, validator.Column)
			}
		}

		if _, err := regexp.Compile(validator.Pattern); err != nil {
			return fmt.Errorf("%s: invalid variable regexp pattern: %w", baseErr, err)
		}
	}

	if v.If != "" {
		if _, err := expr.Compile(v.If); err != nil {
			return fmt.Errorf("invalid 'if' expression: %w", err)
		}
	}

	return nil
}

func (r *VariableValidator) CreateValidatorFunc() func(input string) error {
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
	}
}
