package survey_test

import (
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	"github.com/futurice/jalapeno/pkg/ui/util"
)

func TestPromptUserForValues(t *testing.T) {
	testCases := []struct {
		name                 string
		variables            []recipe.Variable
		existingValues       recipe.VariableValues
		expectedValues       recipe.VariableValues
		input                string
		expectedOutputRegexp string
	}{
		{
			name: "string_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1"},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": "foo",
			},
			input: "foo\r",
		},
		{
			name: "boolean_variable_with_arrows",
			variables: []recipe.Variable{
				{Name: "VAR_1", Confirm: true},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": true,
			},
			input: "→\r",
		},
		{
			name: "boolean_variable_with_keys",
			variables: []recipe.Variable{
				{Name: "VAR_1", Confirm: true},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": true,
			},
			input: "y\r",
		},
		{
			name: "select_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1", Options: []string{"a", "b", "c"}},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": "c",
			},
			input: "↓↓\r",
		},
		{
			name: "multi_select_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1", Options: []string{"a", "b", "c"}, Multi: true},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": []string{"b", "c"},
			},
			input: "↓ ↓ \r",
		},
		{
			name: "table_variable_with_arrows",
			variables: []recipe.Variable{
				{Name: "VAR_1", Columns: []string{"column_1", "column_2"}},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": recipe.TableValue{
					Columns: []string{"column_1", "column_2"},
					Rows:    [][]string{{"foo", "bar"}, {"", "quz"}},
				},
			},
			input: "foo→bar↓quz\r",
		},
		{
			name: "table_variable_with_tabs",
			variables: []recipe.Variable{
				{Name: "VAR_1", Columns: []string{"column_1", "column_2"}},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": recipe.TableValue{
					Columns: []string{"column_1", "column_2"},
					Rows:    [][]string{{"foo", "bar"}, {"baz", "quz"}},
				},
			},
			input: "foo\tbar\tbaz\tquz\r",
		},
		{
			name: "multiple_variables",
			variables: []recipe.Variable{
				{Name: "VAR_1"},
				{Name: "VAR_2", Confirm: true},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": "foo",
				"VAR_2": true,
			},
			input: "foo\ry\r",
		},
		{
			name: "optional_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1", Optional: true},
			},
			expectedValues: recipe.VariableValues{
				"VAR_1": "",
			},
			input:                "\r",
			expectedOutputRegexp: "VAR_1: empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tm := teatest.NewTestModel(
				tt,
				survey.NewModel(tc.variables, tc.existingValues),
				teatest.WithInitialTermSize(300, 100),
			)

			for _, r := range tc.input {
				tm.Send(util.MapUTFRuneToKey(r))
			}

			m := tm.FinalModel(tt, teatest.WithFinalTimeout(time.Second)).(survey.SurveyModel)

			// Assert that the result is correct
			result := m.Values()
			if !reflect.DeepEqual(result, tc.expectedValues) {
				t.Errorf("Unexpected result. Got %v, expected %v", result, tc.expectedValues)
			}

			// Assert that the output is correct
			if tc.expectedOutputRegexp != "" {
				buf := new(strings.Builder)
				_, err := io.Copy(buf, tm.FinalOutput(tt))
				if err != nil {
					t.Fatalf("Could not read output of the survey: %v", err)
				}

				reg := regexp.MustCompile(tc.expectedOutputRegexp)
				if !reg.MatchString(buf.String()) {
					t.Errorf("The output did not match the regular expression \"%s\". Output is: %v", tc.expectedOutputRegexp, buf.String())
				}
			}
		})
	}
}
