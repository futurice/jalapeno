package prompt

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
	"github.com/muesli/reflow/wordwrap"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
)

type SelectModel struct {
	variable        recipe.Variable
	list            list.Model
	styles          style.Styles
	value           string
	showDescription bool
	submitted       bool
	width           int
}

var _ Model = SelectModel{}

type selectItem string

var _ list.Item = selectItem("")

func (i selectItem) FilterValue() string { return "" }

type selectItemDelegate struct{}

var _ list.ItemDelegate = selectItemDelegate{}

func (d selectItemDelegate) Height() int                             { return 1 }
func (d selectItemDelegate) Spacing() int                            { return 0 }
func (d selectItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d selectItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(selectItem)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(string(i))) //nolint:errcheck
}

func NewSelectModel(v recipe.Variable, styles style.Styles) SelectModel {
	items := make([]list.Item, len(v.Options))
	for i := range v.Options {
		items[i] = selectItem(v.Options[i])
	}

	const (
		defaultWidth  = 20
		defaultHeight = 14
	)

	l := list.New(items, selectItemDelegate{}, defaultWidth, defaultHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
			m.value = string(m.list.SelectedItem().(selectItem))
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
		m.list.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectModel) View() string {
	var s strings.Builder
	s.WriteString(m.styles.VariableName.Render(m.variable.Name))
	if m.submitted {
		s.WriteString(fmt.Sprintf(": %s", m.value))
		return s.String()
	}

	if m.variable.Description != "" && !m.showDescription {
		s.WriteString(m.styles.HelpText.Render(" [type ? for more info]"))
	}

	s.WriteRune('\n')
	if m.showDescription {
		s.WriteString(wordwrap.String(m.variable.Description, m.width))
		s.WriteRune('\n')
	}

	s.WriteString(m.list.View())
	return s.String()
}

func (m SelectModel) Name() string {
	return m.variable.Name
}

func (m SelectModel) Value() any {
	return m.value
}

func (m SelectModel) IsSubmitted() bool {
	return m.submitted
}
