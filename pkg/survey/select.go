package survey

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type SelectPromptModel struct {
	variable recipe.Variable
}

var _ PromptModel = SelectPromptModel{}

func NewSelectPromptModel(v recipe.Variable) SelectPromptModel {
	return SelectPromptModel{
		variable: v,
	}
}

func (m SelectPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SelectPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m SelectPromptModel) View() string {
	return ""
}

func (m SelectPromptModel) Value() interface{} {
	return ""
}
