package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/survey/editable"
)

type TableModel struct {
	variable recipe.Variable
	table    editable.Model
}

var _ Model = TableModel{}

func NewTableModel(v recipe.Variable) TableModel {
	cols := make([]editable.Column, len(v.Columns))

	validators := make(map[string][]func(string) error)
	for i, validator := range v.Validators {
		if validator.Column != "" {
			if validators[validator.Column] == nil {
				validators[validator.Column] = make([]func(string) error, 0)
			}

			validators[validator.Column] = append(validators[validator.Column], v.Validators[i].CreateValidatorFunc())
		}
	}

	for i, c := range v.Columns {
		cols[i] = editable.Column{
			Title:      c,
			Width:      len(c),
			Validators: validators[c],
		}
	}
	table := editable.New(editable.WithColumns(cols))
	table.Focus()

	return TableModel{
		variable: v,
		table:    table,
	}
}

func (m TableModel) Init() tea.Cmd {
	return m.table.Init()
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			fmt.Println(recipeutil.RowsToTable(m.variable.Columns, m.table.Values()))
		}
	}
	tm, cmd := m.table.Update(msg)
	m.table = tm.(editable.Model)
	return m, cmd
}

func (m TableModel) View() string {
	return m.table.View()
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
