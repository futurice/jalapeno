package prompt

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/survey/util"
)

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(0)
)

type SelectModel struct {
	variable        recipe.Variable
	list            list.Model
	styles          util.Styles
	value           string
	showDescription bool
	submitted       bool
}

var _ Model = SelectModel{}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(string(i)))
}

func NewSelectModel(v recipe.Variable, styles util.Styles) SelectModel {
	items := make([]list.Item, len(v.Options))
	for i := range v.Options {
		items[i] = item(v.Options[i])
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = paginationStyle

	return SelectModel{
		variable: v,
		list:     l,
		styles:   styles,
	}
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

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
			m.value = string(m.list.SelectedItem().(item))
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectModel) View() (s string) {
	s += m.styles.VariableName.Render(m.variable.Name)
	if m.submitted {
		s += fmt.Sprintf(": %s", m.value)
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

	s += m.list.View()
	return
}

func (m SelectModel) Name() string {
	return m.variable.Name
}

func (m SelectModel) Value() interface{} {
	return m.value
}

func (m SelectModel) IsSubmitted() bool {
	return m.submitted
}
