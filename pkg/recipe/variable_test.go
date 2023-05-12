package recipe

import "testing"

func TestVariableRegExpValidation(t *testing.T) {
	variable := &Variable{
		Name:        "foo",
		Description: "foo description",
		RegExp: VariableRegExpValidator{
			Pattern: "^[a-zA-Z0-9_.()-]{0,89}[a-zA-Z0-9_()-]$",
		},
	}

	validatorFunc := variable.RegExp.CreateValidatorFunc()

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
