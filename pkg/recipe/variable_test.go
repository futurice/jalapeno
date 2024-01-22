package recipe

import (
	"strings"
	"testing"
)

func TestVariableValidation(t *testing.T) {
	scenarios := []struct {
		name        string
		variable    Variable
		expectedErr string
	}{
		{
			"empty name",
			Variable{},
			"variable name is required",
		},
		{
			"variable name starts with number",
			Variable{
				Name: "1foo",
			},
			"variable name can not start with a number",
		},
		{
			"both `confirm` and `options` defined",
			Variable{
				Name:    "foo",
				Confirm: true,
				Options: []string{"foo", "bar"},
			},
			"`confirm` and `options` properties can not be defined",
		},
		{
			"both `confirm` and `columns` defined",
			Variable{
				Name:    "foo",
				Confirm: true,
				Columns: []string{"foo", "bar"},
			},
			"`confirm` and `columns` properties can not be defined",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := scenario.variable.Validate()
			if err != nil {
				if !strings.Contains(err.Error(), scenario.expectedErr) {
					t.Errorf("Expected error '%s', got '%s'", scenario.expectedErr, err.Error())
				}
			} else {
				if scenario.expectedErr != "" {
					t.Errorf("Expected error '%s', got nil", scenario.expectedErr)
				}
			}
		})
	}
}

func TestVariableRegExpValidation(t *testing.T) {
	variable := &Variable{
		Name:        "foo",
		Description: "foo description",
		Validators: []VariableValidator{
			{
				Pattern: "^[a-zA-Z0-9_.()-]{0,89}[a-zA-Z0-9_()-]$",
			},
		},
	}

	validatorFunc := variable.Validators[0].CreateValidatorFunc()

	err := validatorFunc("")
	if err == nil {
		t.Error("Incorrectly validated empty string")
	}

	err = validatorFunc("foo bar baz")
	if err == nil {
		t.Error("Incorrectly validated string with whitespace")
	}

	err = validatorFunc("012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789x")
	if err == nil {
		t.Error("Incorrectly validated too long string")
	}

	err = validatorFunc("valid-except-for-#¤%&?")
	if err == nil {
		t.Error("Incorrectly validated special characters")
	}

	err = validatorFunc("valid-except.")
	if err == nil {
		t.Error("Incorrectly validated string ending in a period")
	}

	err = validatorFunc("valid-all-the-way-to-11")
	if err != nil {
		t.Error("Incorrectly invalidated valid string")
	}
}
