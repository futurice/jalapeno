package prompt

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/editable"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
	"github.com/muesli/reflow/wordwrap"
)

type TableModel struct {
	variable        recipe.Variable
	table           editable.Model
	styles          style.Styles
	submitted       bool
	showDescription bool
	width           int

	// Save the table as CSV for the final output. This speeds up the
	// rendering after the user has submitted the form.
	tableAsCSV string
	err        error
}

var _ Model = TableModel{}

func NewTableModel(v recipe.Variable, styles style.Styles) TableModel {
	cols := make([]editable.Column, len(v.Columns))

	validators := make(map[string][]func([]string, [][]string, string) error)
	for i, validator := range v.Validators {
		if validator.Column != "" {
			if validators[validator.Column] == nil {
				validators[validator.Column] = make([]func([]string, [][]string, string) error, 0)
			}

			if validator.Pattern != "" {
				regexValidator, err := v.Validators[i].CreateValidatorFunc()
				if err == nil {
					validators[validator.Column] = append(validators[validator.Column],
						func(cols []string, rows [][]string, input string) error {
							return regexValidator(input)
						})
				}
			} else {
				validatorFn, err := validator.CreateTableValidatorFunc()
				if err == nil {
					validators[validator.Column] = append(validators[validator.Column], validatorFn)
				}
			}
		}
	}
	for i, c := range v.Columns {
		cols[i] = editable.Column{
			Title:      c,
			Width:      len(c),
			Validators: validators[c],
		}
	}

	table := editable.NewModel(editable.WithColumns(cols), editable.IsOptional(v.Optional))
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
		switch msg.Type {
		case tea.KeyEnter:
			// Validate the table. If there are errors, don't submit the form.
			m.table.Validate()
			if errs := m.table.Errors(); len(errs) != 0 {
				return m, nil
			}

			if !m.variable.Optional && m.IsEmpty() {
				m.err = errors.New("table can not be empty since the variable is not optional")
				return m, nil
			}

			m.submitted = true
			m.tableAsCSV = m.ValueAsCSV()
			m.table.Blur()
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "?":
				if m.variable.Description != "" && !m.showDescription {
					m.showDescription = true
					return m, nil
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	tm, cmd := m.table.Update(msg)
	m.table = tm.(editable.Model)
	return m, cmd
}

func (m TableModel) View() string {
	var s strings.Builder
	s.WriteString(m.styles.VariableName.Render(m.variable.Name))

	if m.submitted {
		s.WriteString(": ")

		if m.IsEmpty() {
			s.WriteString(m.styles.HelpText.Render("empty"))
		} else {
			s.WriteString(m.tableAsCSV)
		}

		s.WriteString(m.tableAsCSV)
		return s.String()
	}

	if m.variable.Optional {
		s.WriteString(m.styles.HelpText.Render(" (optional)"))
	}

	if !m.showDescription {
		s.WriteString(m.styles.HelpText.Render(" [type ? for more info]"))
	}

	s.WriteRune('\n')
	if m.showDescription {
		if m.variable.Description != "" {
			s.WriteString(wordwrap.String(m.variable.Description, m.width))
			s.WriteRune('\n')
		}
		s.WriteString(wordwrap.String(m.styles.HelpText.Render(`Table controls:
- arrow keys: to move between cells
- tab: to move to the next cells
- ctrl+n or move past last row: create a new row 
`), m.width))
		s.WriteRune('\n')
	}

	s.WriteString(m.table.View())

	if m.err != nil {
		errMsg := m.err.Error()
		errMsg = strings.ToUpper(errMsg[:1]) + errMsg[1:]
		s.WriteString(m.styles.ErrorText.Render(errMsg))
	}

	return s.String()
}

func (m TableModel) Name() string {
	return m.variable.Name
}

func (m TableModel) Value() interface{} {
	return recipe.TableValue{
		Columns: m.variable.Columns,
		Rows:    m.table.Values(),
	}
}

func (m TableModel) IsSubmitted() bool {
	return m.submitted
}

var (
	csvSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")).SetString(",")
	csvNewLine   = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")).SetString("\\n")
)

func (m TableModel) ValueAsCSV() string {
	rows := m.table.Values()
	s := ""
	for y := range rows {
		s += strings.Join(rows[y], csvSeparator.String())
		if y < len(rows)-1 {
			s += csvNewLine.String()
		}
	}

	return s
}

func (m TableModel) IsEmpty() bool {
	return len(m.table.Values()) == 0
}
