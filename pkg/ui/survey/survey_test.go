package survey_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey"
	"github.com/futurice/jalapeno/pkg/ui/util"
)

func TestPromptUserForValues(t *testing.T) {
	testCases := []struct {
		name           string
		variables      []recipe.Variable
		existingValues recipe.VariableValues
		expected       recipe.VariableValues
		input          string
	}{
		{
			name: "string_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1"},
			},
			expected: recipe.VariableValues{
				"VAR_1": "foo",
			},
			input: "foo\r",
		},
		{
			name: "select_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1", Options: []string{"a", "b", "c"}},
			},
			expected: recipe.VariableValues{
				"VAR_1": "c",
			},
			input: "↓↓\r",
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
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Unexpected result. Got %v, expected %v", result, tc.expected)
			}
		})
	}
}
