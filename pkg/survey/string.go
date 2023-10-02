package survey

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type StringPromptModel struct {
	variable  recipe.Variable
	textInput textinput.Model
	styles    Styles
	err       error
}

var _ PromptModel = StringPromptModel{}

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

func NewStringPromptModel(v recipe.Variable) StringPromptModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	if v.Default != "" {
		ti.SetValue(v.Default)
	}

	return StringPromptModel{
		variable:  v,
		textInput: ti,
		err:       nil,
		styles:    DefaultStyles(),
	}
}

func (m StringPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StringPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case FocusMsg:
		m.textInput.Focus()
		return m, nil
	case BlurMsg:
		m.textInput.Blur()
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StringPromptModel) View() (s string) {
	s += fmt.Sprintf("%s:\n", m.styles.VariableName.Render(m.variable.Name))

	if m.textInput.Focused() {
		s += m.variable.Description
		s += "\n"
	}

	s += m.textInput.View()

	return
}

func (m StringPromptModel) Value() interface{} {
	return m.textInput.Value()
}
