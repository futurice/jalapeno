package recipeutil_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
)

func TestParsePredefinedValues(t *testing.T) {
	tests := []struct {
		name        string
		vars        []recipe.Variable
		envs        [][2]string
		flags       []string
		expected    recipe.VariableValues
		expectedErr error
	}{
		{
			name:  "pass_flag",
			vars:  []recipe.Variable{{Name: "test_var"}},
			flags: []string{"test_var=value"},
			expected: recipe.VariableValues{
				"test_var": "value",
			},
		},
		{
			name: "pass_env",
			vars: []recipe.Variable{{Name: "test_var"}},
			envs: [][2]string{{"test_var", "value"}},
			expected: recipe.VariableValues{
				"test_var": "value",
			},
		},
		{
			name:        "unknown_var_flag",
			vars:        []recipe.Variable{},
			flags:       []string{"test_var=value"},
			expectedErr: recipeutil.ErrVarNotDefinedInRecipe,
		},
		{
			name:        "unknown_var_env",
			vars:        []recipe.Variable{},
			envs:        [][2]string{{"test_var", "value"}},
			expectedErr: recipeutil.ErrVarNotDefinedInRecipe,
		},
		{
			name:  "flag_overrides_env",
			vars:  []recipe.Variable{{Name: "test_var"}},
			envs:  [][2]string{{"test_var", "first"}},
			flags: []string{"test_var=second"},
			expected: recipe.VariableValues{
				"test_var": "second",
			},
		},
		{
			name:  "multiple_vars",
			vars:  []recipe.Variable{{Name: "test_var_1"}, {Name: "test_var_2"}},
			flags: []string{"test_var_1=value1"},
			envs:  [][2]string{{"test_var_2", "value2"}},
			expected: recipe.VariableValues{
				"test_var_1": "value1",
				"test_var_2": "value2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, env := range test.envs {
				envName := recipeutil.ValueEnvVarPrefix + env[0]
				err := os.Setenv(envName, env[1])
				if err != nil {
					t.Fatal("failed to set environment variable")
				}
				defer os.Unsetenv(envName)
			}

			actual, err := recipeutil.ParseProvidedValues(test.vars, test.flags)
			if err != nil {
				if test.expectedErr == nil {
					t.Fatalf("parser returned error when not expected, error: %+v", err)
				}

				if !errors.Is(err, test.expectedErr) {
					t.Fatalf("parser error did not match expected error: expected: %s, actual: %s", test.expectedErr, err)
				}
			} else if test.expectedErr != nil {
				t.Fatal("parser did not return error when expected")
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("parsed values had non-expected value, expected %+v, actual: %+v", test.expected, actual)
			}
		})
	}
}

func TestMergeValues(t *testing.T) {
	tests := []struct {
		name     string
		a        recipe.VariableValues
		b        recipe.VariableValues
		expected recipe.VariableValues
	}{
		{
			name: "merge",
			a: recipe.VariableValues{
				"FOO": "foo",
			},
			b: recipe.VariableValues{
				"BAR": "bar",
			},
			expected: recipe.VariableValues{
				"FOO": "foo",
				"BAR": "bar",
			},
		},
		{
			name: "overlap",
			a: recipe.VariableValues{
				"FOO": "foo",
			},
			b: recipe.VariableValues{
				"FOO": "bar",
			},
			expected: recipe.VariableValues{
				"FOO": "bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := recipeutil.MergeValues(test.a, test.b)

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("Merged values had non-expected value, expected %+v, actual: %+v", test.expected, actual)
			}
		})
	}

}

func TestFilterVariables(t *testing.T) {
	tests := []struct {
		name      string
		variables []recipe.Variable
		values    recipe.VariableValues
		expected  []recipe.Variable
	}{
		{
			name: "value_exists",
			variables: []recipe.Variable{
				{Name: "FOO"},
			},
			values: recipe.VariableValues{
				"FOO": "foo",
			},
			expected: []recipe.Variable{},
		},
		{
			name: "missing_value",
			variables: []recipe.Variable{
				{Name: "FOO"},
			},
			values: recipe.VariableValues{
				"BAR": "bar",
			},
			expected: []recipe.Variable{
				{Name: "FOO"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := recipeutil.FilterVariablesWithoutValues(test.variables, test.values)

			if !reflect.DeepEqual(test.expected, actual) {
				t.Fatalf("Merged values had non-expected value, expected %+v, actual: %+v", test.expected, actual)
			}
		})
	}
}
