package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type SelectModel struct {
	variable recipe.Variable
}

var _ Model = SelectModel{}

func NewSelectModel(v recipe.Variable) SelectModel {
	return SelectModel{
		variable: v,
	}
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m SelectModel) View() string {
	return ""
}

func (m SelectModel) Value() interface{} {
	return ""
}

func (m SelectModel) IsSubmitted() bool {
	return false
}
