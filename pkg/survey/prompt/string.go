package prompt

import (
	"fmt"
	"strings"

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
	ErrorText    lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		VariableName: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")),
		ErrorText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")),
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
			if err := m.Validate(); err != nil {
				m.err = err
				return m, nil
			}
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

	if m.textInput.Focused() && m.err != nil {
		s += "\n"
		errMsg := m.err.Error()
		errMsg = strings.ToUpper(errMsg[:1]) + errMsg[1:]
		s += m.styles.ErrorText.Render(errMsg)
	}

	return
}

func (m StringModel) Name() string {
	return m.variable.Name
}

func (m StringModel) Value() interface{} {
	return m.textInput.Value()
}

func (m StringModel) IsSubmitted() bool {
	return m.submitted
}

func (m StringModel) Validate() error {
	if !m.variable.Optional && m.textInput.Value() == "" {
		return util.ErrRequired
	}

	if m.variable.RegExp.Pattern != "" {
		validator := m.variable.RegExp.CreateValidatorFunc()
		if err := validator(m.textInput.Value()); err != nil {
			return fmt.Errorf("%w: %s", util.ErrRegExFailed, err)
		}
	}

	return nil
}
