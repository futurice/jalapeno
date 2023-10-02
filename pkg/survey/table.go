package survey

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type TablePromptModel struct {
	variable recipe.Variable
}

var _ PromptModel = TablePromptModel{}

func NewTablePromptModel(v recipe.Variable) TablePromptModel {
	return TablePromptModel{
		variable: v,
	}
}

func (m TablePromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TablePromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TablePromptModel) View() string {
	return ""
}

func (m TablePromptModel) Value() interface{} {
	return ""
}
