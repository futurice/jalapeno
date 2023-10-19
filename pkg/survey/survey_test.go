package survey_test

import (
	"reflect"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/survey"
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
			input: "foo\n",
		},
		{
			name: "select_variable",
			variables: []recipe.Variable{
				{Name: "VAR_1", Options: []string{"a", "b", "c"}},
			},
			expected: recipe.VariableValues{
				"VAR_1": "c",
			},
			input: "↓↓\n",
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
				tm.Send(RuneToKey(r))
			}

			m := tm.FinalModel(tt, teatest.WithFinalTimeout(time.Second)).(survey.SurveyModel)
			m.Values()

			// Assert that the result is correct
			result := m.Values()
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Unexpected result. Got %v, expected %v", result, tc.expected)
			}
		})
	}
}

func RuneToKey(r rune) tea.KeyMsg {
	switch r {
	case '\n':
		return tea.KeyMsg{
			Type: tea.KeyEnter,
		}
	case '↑':
		return tea.KeyMsg{
			Type: tea.KeyUp,
		}
	case '↓':
		return tea.KeyMsg{
			Type: tea.KeyDown,
		}
	case '←':
		return tea.KeyMsg{
			Type: tea.KeyLeft,
		}
	case '→':
		return tea.KeyMsg{
			Type: tea.KeyRight,
		}
	default:
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
		}
	}
}
