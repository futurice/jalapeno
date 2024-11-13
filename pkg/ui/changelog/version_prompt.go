package changelog

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
)

type VersionModel struct {
	list  list.Model
	value string
	width int
}

var _ tea.Model = VersionModel{}

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

	fmt.Fprint(w, fn(string(i)))
}

func NewSelectModel(options []string) VersionModel {
	items := make([]list.Item, len(options))
	for i := range options {
		items[i] = selectItem(options[i])
	}

	l := list.New(items, selectItemDelegate{}, 6, 5)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	return VersionModel{
		list: l,
	}
}

func (m VersionModel) Init() tea.Cmd {
	return nil
}

func (m VersionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.value = string(m.list.SelectedItem().(selectItem))
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.list.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m VersionModel) View() string {
	var s strings.Builder

	s.WriteString(wordwrap.String("Select which version to bump:", m.width))
	s.WriteRune('\n')
	s.WriteString(m.list.View())

	return s.String()
}

func (m VersionModel) Value() string {
	return m.value
}
