package prompt

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	tea.Model
	IsSubmitted() bool
	Value() interface{}
}

var _ tea.Model = Model(nil)
