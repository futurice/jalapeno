package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type ConfirmModel struct {
	variable recipe.Variable
}

var _ Model = ConfirmModel{}

func NewConfirmModel(v recipe.Variable) ConfirmModel {
	return ConfirmModel{
		variable: v,
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ConfirmModel) View() string {
	return ""
}

func (m ConfirmModel) Name() string {
	return m.variable.Name
}

func (m ConfirmModel) Value() interface{} {
	return true
}

func (m ConfirmModel) IsSubmitted() bool {
	return false
}
