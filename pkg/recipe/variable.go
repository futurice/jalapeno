package recipe

import (
	"errors"
	"fmt"
	"regexp"
)

type Variable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`

	// Default value for the variable
	Default string `yaml:"default,omitempty"`

	// If set to true, the prompt will be yes/no question, and the value type will be boolean
	Confirm bool `yaml:"confirm,omitempty"`

	// If set to true, the variable can be left empty
	Optional bool `yaml:"optional,omitempty"`

	// The user selects the value from a list of options
	Options []string `yaml:"options,omitempty"`

	// Regular expression validator for the variable value
	RegExp VariableRegExpValidator `yaml:"regexp,omitempty"`
}

type VariableRegExpValidator struct {
	Pattern string `yaml:"pattern,omitempty"`

	// If the RegExp validation fails, this help message will be showed to the user
	Help string `yaml:"help,omitempty"`
}

type VariableValues map[string]interface{}

func (v *Variable) Validate() error {
	// TODO
	return nil
}

func (r *VariableRegExpValidator) CreateValidatorFunc() (func(input interface{}) error, error) {
	reg, err := regexp.Compile(r.Pattern)
	if err != nil {
		return nil, err
	}

	validator := func(input interface{}) error {
		if match := reg.MatchString(fmt.Sprint(input)); !match {
			if r.Help != "" {
				return errors.New(r.Help)
			} else {
				return errors.New("the input did not match the regexp pattern")
			}
		}
		return nil
	}

	return validator, nil
}
