package survey

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type ConfirmPromptModel struct {
	variable recipe.Variable
}

var _ PromptModel = ConfirmPromptModel{}

func NewConfirmPromptModel(v recipe.Variable) ConfirmPromptModel {
	return ConfirmPromptModel{
		variable: v,
	}
}

func (m ConfirmPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConfirmPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ConfirmPromptModel) View() string {
	return ""
}

func (m ConfirmPromptModel) Value() interface{} {
	return true
}
