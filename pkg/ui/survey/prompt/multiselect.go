package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
	"github.com/muesli/reflow/wordwrap"
)

var (
	cursorStyle                  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	multiSelectItemStyle         = lipgloss.NewStyle()
	selectedMultiSelectItemStyle = lipgloss.NewStyle().Inherit(multiSelectItemStyle).Foreground(lipgloss.Color("205"))
)

type MultiSelectModel struct {
	variable        recipe.Variable
	styles          style.Styles
	items           []multiSelectItem
	index           int
	showDescription bool
	submitted       bool
	width           int
}

var _ Model = MultiSelectModel{}

type multiSelectItem struct {
	value   string
	checked bool
}

func NewMultiSelectModel(v recipe.Variable, styles style.Styles) MultiSelectModel {
	items := make([]multiSelectItem, len(v.Options))
	for i := range v.Options {
		items[i] = multiSelectItem{
			value:   v.Options[i],
			checked: false,
		}
	}

	return MultiSelectModel{
		variable: v,
		styles:   styles,
		items:    items,
	}
}

func (m MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (m MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
		case tea.KeySpace:
			m.items[m.index].checked = !m.items[m.index].checked
		case tea.KeyDown:
			if m.index < len(m.items)-1 {
				m.index++
			}
		case tea.KeyUp:
			if m.index > 0 {
				m.index--
			}
		case tea.KeyHome:
			m.index = 0
		case tea.KeyEnd:
			m.index = len(m.items) - 1
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "?":
				if !m.showDescription {
					m.showDescription = true
					return m, nil
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, cmd
}

func (m MultiSelectModel) View() string {
	var s strings.Builder
	s.WriteString(m.styles.VariableName.Render(m.variable.Name))

	if m.submitted {
		s.WriteString(": ")

		if values := m.getSelectedValues(); len(values) == 0 {
			s.WriteString(m.styles.HelpText.Render("empty"))
		} else {
			s.WriteString(strings.Join(values, ", "))
		}

		return s.String()
	}

	if m.variable.Optional {
		s.WriteString(m.styles.HelpText.Render(" (optional)"))
	}

	if m.variable.Description != "" && !m.showDescription {
		s.WriteString(m.styles.HelpText.Render(" [type ? for more info]"))
	}

	s.WriteRune('\n')
	if m.showDescription {
		if m.variable.Description != "" {
			s.WriteString(wordwrap.String(m.variable.Description, m.width))
			s.WriteRune('\n')
		}
		s.WriteString(wordwrap.String(m.styles.HelpText.Render(`Controls:
- space: select the item
- up/down arrow keys: to move between items
- home: to move to the first item
- end: to move to the last item
`), m.width))
		s.WriteRune('\n')
	}

	for i, item := range m.items {
		if i == m.index {
			s.WriteString(fmt.Sprintf("%s ", cursorStyle.Render("❯")))
		} else {
			s.WriteString("  ")
		}

		if item.checked {
			s.WriteString(selectedMultiSelectItemStyle.Render(fmt.Sprintf("%s ✔", item.value)))
		} else {
			s.WriteString(multiSelectItemStyle.Render(item.value))
		}

		if i < len(m.items)-1 {
			s.WriteRune('\n')
		}
	}

	return s.String()
}

func (m MultiSelectModel) getSelectedValues() []string {
	selectedValues := make([]string, 0, len(m.items))
	for _, item := range m.items {
		if item.checked {
			selectedValues = append(selectedValues, item.value)
		}
	}

	return selectedValues
}

func (m MultiSelectModel) Name() string {
	return m.variable.Name
}

func (m MultiSelectModel) Value() interface{} {
	return m.getSelectedValues()
}

func (m MultiSelectModel) IsSubmitted() bool {
	return m.submitted
}
