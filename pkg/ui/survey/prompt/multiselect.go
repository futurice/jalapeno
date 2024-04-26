package prompt

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
)

type MultiSelectModel struct {
	variable        recipe.Variable
	styles          style.Styles
	value           recipe.MultiSelectValue
	showDescription bool
	submitted       bool
	width           int
}

var _ Model = MultiSelectModel{}

type multiSelectItem string

var _ list.Item = multiSelectItem("")

func (i multiSelectItem) FilterValue() string { return "" }

type itemDelegate struct{}

var _ list.ItemDelegate = itemDelegate{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(multiSelectItem)
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

func NewMultiSelectModel(v recipe.Variable, styles style.Styles) MultiSelectModel {
	items := make([]list.Item, len(v.Options))
	for i := range v.Options {
		items[i] = multiSelectItem(v.Options[i])
	}

	return MultiSelectModel{
		variable: v,
		styles:   styles,
	}
}

func (m MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (m MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m MultiSelectModel) View() string {
	var s strings.Builder
	return s.String()
}

func (m MultiSelectModel) Name() string {
	return m.variable.Name
}

func (m MultiSelectModel) Value() interface{} {
	return m.value
}

func (m MultiSelectModel) IsSubmitted() bool {
	return m.submitted
}
