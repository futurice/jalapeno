package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/survey/util"
)

type ConfirmModel struct {
	variable        recipe.Variable
	styles          util.Styles
	value           bool
	submitted       bool
	showDescription bool
}

var _ Model = ConfirmModel{}

func NewConfirmModel(v recipe.Variable, styles util.Styles) ConfirmModel {
	return ConfirmModel{
		variable: v,
		styles:   styles,
		value:    v.Default == "true",
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			if m.variable.Description != "" && !m.showDescription {
				m.showDescription = true
				return m, nil
			}
		case "y", "Y":
			m.value = true
		case "n", "N":
			m.value = false
		}
		switch msg.Type {
		case tea.KeyRight:
			m.value = true
		case tea.KeyLeft:
			m.value = false
		case tea.KeyEnter:
			m.submitted = true
		}
	}

	return m, nil
}

func (m ConfirmModel) View() (s string) {
	s += m.styles.VariableName.Render(m.variable.Name)
	if m.submitted {
		s += ": "
		if m.value {
			s += "Yes"
		} else {
			s += "No"
		}
		return
	}

	if m.variable.Description != "" && !m.showDescription {
		s += m.styles.HelpText.Render(" [type ? for more info]")
	}

	s += "\n"
	if m.showDescription {
		s += m.variable.Description
		s += "\n"
	}

	if m.value {
		s += fmt.Sprintf("> No/%s", m.styles.Bold.Render("Yes"))
	} else {
		s += fmt.Sprintf("> %s/Yes", m.styles.Bold.Render("No"))
	}

	return
}

func (m ConfirmModel) Name() string {
	return m.variable.Name
}

func (m ConfirmModel) Value() interface{} {
	return m.value
}

func (m ConfirmModel) IsSubmitted() bool {
	return m.submitted
}
