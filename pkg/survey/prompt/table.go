package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type TableModel struct {
	variable recipe.Variable
}

var _ Model = TableModel{}

func NewTableModel(v recipe.Variable) TableModel {
	return TableModel{
		variable: v,
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TableModel) View() string {
	return ""
}

func (m TableModel) Name() string {
	return m.variable.Name
}

func (m TableModel) Value() interface{} {
	return ""
}

func (m TableModel) IsSubmitted() bool {
	return false
}
