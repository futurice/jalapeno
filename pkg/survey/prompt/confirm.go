package prompt

import (
	"fmt"
	"strings"

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
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
		case tea.KeyRight:
			m.value = true
		case tea.KeyLeft:
			m.value = false
		case tea.KeyRunes:
			switch string(msg.Runes) {
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
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	var s strings.Builder
	s.WriteString(m.styles.VariableName.Render(m.variable.Name))
	if m.submitted {
		s.WriteString(": ")
		if m.value {
			s.WriteString("Yes")
		} else {
			s.WriteString("No")
		}
		return s.String()
	}

	if m.variable.Description != "" && !m.showDescription {
		s.WriteString(m.styles.HelpText.Render(" [type ? for more info]"))
	}

	s.WriteRune('\n')
	if m.showDescription {
		s.WriteString(m.variable.Description)
		s.WriteRune('\n')
	}

	if m.value {
		s.WriteString(fmt.Sprintf("> No/%s", m.styles.Bold.Render("Yes")))
	} else {
		s.WriteString(fmt.Sprintf("> %s/Yes", m.styles.Bold.Render("No")))
	}

	return s.String()
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
