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
		{
			"both `options` and `columns` defined",
			Variable{
				Name:    "foo",
				Options: []string{"foo", "bar"},
				Columns: []string{"foo", "bar"},
			},
			"`options` and `columns` properties can not be defined",
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
		Validators: []VariableValidator{
			{
				Pattern: "^[a-zA-Z0-9_.()-]{0,89}[a-zA-Z0-9_()-]$",
			},
		},
	}

	validatorFunc, err := variable.Validators[0].CreateValidatorFunc()
	if err != nil {
		t.Error("Validator function creation failed")
	}

	err = validatorFunc("")
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

func TestUniqueColumnValidation(t *testing.T) {
	variable := &Variable{
		Validators: []VariableValidator{
			{
				Unique: true,
				Column: "COL_1",
			},
		},
	}

	validatorFunc, err := variable.Validators[0].CreateTableValidatorFunc()
	if err != nil {
		t.Error("Validator function creation failed")
	}

	cols := []string{"COL_1", "COL_2"}

	err = validatorFunc(
		cols,
		[][]string{
			{"0_0", "0_1"},
			{"1_0", "1_1"},
			{"2_0", "2_1"},
		},
		"")
	if err != nil {
		t.Error("Incorrectly invalidated valid data")
	}

	err = validatorFunc(
		cols,
		[][]string{
			{"0_0", "0_1"},
			{"0_0", "1_1"},
			{"2_0", "2_1"},
		},
		"")
	if err == nil {
		t.Error("Incorrectly validated invalid data")
	}

	err = validatorFunc(
		cols,
		[][]string{
			{"0_0", "0_1"},
			{"1_0", "0_1"},
			{"2_0", "0_1"},
		},
		"")
	if err != nil {
		t.Error("Incorrectly invalidated valid data")
	}
}
