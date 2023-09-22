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

	// Regular expression validator for the variable value
	RegExp VariableRegExpValidator `yaml:"regexp,omitempty"`

	// Makes the variable conditional based on the result of the expression. The result of the evaluation needs to be a boolean value. Uses https://github.com/antonmedv/expr
	If string `yaml:"if,omitempty"`
}

type VariableRegExpValidator struct {
	// Regular expression pattern to match the input against
	Pattern string `yaml:"pattern,omitempty"`

	// If the regular expression validation fails, this help message will be shown to the user
	Help string `yaml:"help,omitempty"`
}

// VariableValues stores values for each variable. In practice the value can bee either string or boolean.
type VariableValues map[string]interface{}

func (v *Variable) Validate() error {
	if v.Name == "" {
		return errors.New("variable name is required")
	}
	if v.Confirm && len(v.Options) > 0 {
		return errors.New("`confirm` and `options` properties can not be defined at the same time")
	}
	if v.RegExp.Pattern != "" {
		if _, err := regexp.Compile(v.RegExp.Pattern); err != nil {
			return fmt.Errorf("invalid variable regexp pattern: %w", err)
		}
	}
	if v.If != "" {
		if _, err := expr.Compile(v.If); err != nil {
			return fmt.Errorf("invalid variable 'if' expression: %w", err)
		}
	}
	return nil
}

func (r *VariableRegExpValidator) CreateValidatorFunc() func(input interface{}) error {
	reg := regexp.MustCompile(r.Pattern)

	return func(input interface{}) error {
		if match := reg.MatchString(fmt.Sprint(input)); !match {
			if r.Help != "" {
				return errors.New(r.Help)
			} else {
				return errors.New("the input did not match the regexp pattern")
			}
		}
		return nil
	}
}
