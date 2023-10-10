package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/survey/editable"
	"github.com/futurice/jalapeno/pkg/survey/util"
)

type TableModel struct {
	variable        recipe.Variable
	table           editable.Model
	styles          util.Styles
	submitted       bool
	showDescription bool

	// Save the table as CSV for the final output. This speeds up the
	// rendering when the user has submitted the form.
	tableAsCSV string
}

var _ Model = TableModel{}

func NewTableModel(v recipe.Variable, styles util.Styles) TableModel {
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
	table := editable.NewModel(editable.WithColumns(cols))
	table.Focus()

	return TableModel{
		variable: v,
		table:    table,
		styles:   styles,
	}
}

func (m TableModel) Init() tea.Cmd {
	return m.table.Init()
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			if m.variable.Description != "" && !m.showDescription {
				m.showDescription = true
				return m, nil
			}
		}
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
			m.tableAsCSV = m.ValueAsCSV()
			m.table.Blur()
		}
	}
	tm, cmd := m.table.Update(msg)
	m.table = tm.(editable.Model)
	return m, cmd
}

func (m TableModel) View() (s string) {
	s += m.styles.VariableName.Render(m.variable.Name)

	if m.submitted {
		s += ": "
		s += m.tableAsCSV
		return
	}

	if m.variable.Description != "" && !m.showDescription {
		s += m.styles.HelpText.Render(" [type ? for more info]")
	}

	s += "\n"
	if m.showDescription {
		s += m.variable.Description
	}
	s += "\n"

	s += m.table.View()
	return
}

func (m TableModel) Name() string {
	return m.variable.Name
}

func (m TableModel) Value() interface{} {
	values, _ := recipeutil.RowsToTable(m.variable.Columns, m.table.Values())
	return values
}

func (m TableModel) IsSubmitted() bool {
	return m.submitted
}

func (m TableModel) ValueAsCSV() string {
	rows := m.table.Values()
	s := ""
	for y := range rows {
		for x := range rows[y] {
			s += rows[y][x]
			if x < len(rows[y])-1 {
				s += ","
			}
		}
		if y < len(rows)-1 {
			s += "\\n"
		}
	}

	return s
}
