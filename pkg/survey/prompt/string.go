package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/survey/util"
)

type StringModel struct {
	variable        recipe.Variable
	textInput       textinput.Model
	styles          Styles
	submitted       bool
	showDescription bool
	err             error
}

var _ Model = StringModel{}

type Styles struct {
	VariableName lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		VariableName: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")),
	}
}

func NewStringModel(v recipe.Variable) StringModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	if v.Default != "" {
		ti.SetValue(v.Default)
	}

	return StringModel{
		variable:  v,
		textInput: ti,
		err:       nil,
		styles:    DefaultStyles(),
	}
}

func (m StringModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			if m.textInput.Value() == "" && m.variable.Description != "" && !m.showDescription {
				m.showDescription = true
				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
		}
	case util.FocusMsg:
		m.textInput.Focus()
		m.textInput.Prompt = "> "
		return m, nil
	case util.BlurMsg:
		m.textInput.Blur()
		m.textInput.Prompt = ""
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StringModel) View() (s string) {
	s += m.styles.VariableName.Render(m.variable.Name)

	if m.textInput.Focused() {
		if m.variable.Description != "" && !m.showDescription {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
			s += style.Render(" [type ? for more info]")
		}

		s += "\n"
		if m.showDescription {
			s += m.variable.Description
			s += "\n"
		}
	} else {
		s += ": "
	}

	s += m.textInput.View()

	return
}

func (m StringModel) Value() interface{} {
	return m.textInput.Value()
}

func (m StringModel) IsSubmitted() bool {
	return m.submitted
}
