package util

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	ErrRequired    = errors.New("value can not be empty")
	ErrRegExFailed = errors.New("validation failed")
	ErrUserAborted = errors.New("user aborted")
)

func MapUTFRuneToKey(r rune) tea.KeyMsg {
	switch r {
	case '\r':
		return tea.KeyMsg{
			Type: tea.KeyEnter,
		}
	case '\t':
		return tea.KeyMsg{
			Type: tea.KeyTab,
		}
	case ' ':
		return tea.KeyMsg{
			Type: tea.KeySpace,
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
