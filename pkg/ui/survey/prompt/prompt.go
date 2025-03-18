package prompt

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	tea.Model
	IsSubmitted() bool
	Name() string
	Value() any
}

var _ tea.Model = Model(nil)
