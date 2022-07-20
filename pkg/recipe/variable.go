package recipe

import (
	"errors"
	"fmt"
	"regexp"
)

type Variable struct {
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description,omitempty"`
	Default     string                  `yaml:"default,omitempty"`
	Optional    bool                    `yaml:"optional,omitempty"`
	Options     []string                `yaml:"options,omitempty"`
	RegExp      VariableRegExpValidator `yaml:"regexp,omitempty"`
}

type VariableRegExpValidator struct {
	Pattern string `yaml:"pattern,omitempty"`
	Help    string `yaml:"help,omitempty"`
}

type VariableValues map[string]string

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
