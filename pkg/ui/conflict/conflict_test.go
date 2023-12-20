package conflict_test

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/futurice/jalapeno/pkg/ui/conflict"
)

func TestSolveFileConflict(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		fileA    []byte
		fileB    []byte
		input    string
		expected []byte
	}{
		{
			name:     "no_answer",
			filePath: "README.md",
			fileA:    []byte("foo"),
			fileB:    []byte("bar"),
			input:    "n\n",
			expected: []byte("foo"),
		},
		{
			name:     "yes_answer",
			filePath: "README.md",
			fileA:    []byte("foo"),
			fileB:    []byte("bar"),
			input:    "y\n",
			expected: []byte("bar"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tm := teatest.NewTestModel(
				tt,
				conflict.NewModel(tc.filePath, tc.fileA, tc.fileB),
				teatest.WithInitialTermSize(300, 100),
			)

			for _, r := range tc.input {
				tm.Send(RuneToKey(r))
			}

			m := tm.FinalModel(tt, teatest.WithFinalTimeout(time.Second)).(conflict.Model)

			// Assert that the result is correct
			result := m.Result()
			if !bytes.Equal(result, tc.expected) {
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
