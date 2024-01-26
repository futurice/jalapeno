package editable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	KeyMap KeyMap

	cols     []Column
	rows     []Row
	cursorX  int
	cursorY  int
	focus    bool
	optional bool

	width int

	styles Styles
	table  *table.Table
}

var _ tea.Model = Model{}
var _ table.Data = Model{}

type Row []Cell

type Cell struct {
	input textinput.Model
	err   error
}

type Column struct {
	Title      string
	Width      int
	Validators []func(string) error
}

type KeyMap struct {
	CellUp     key.Binding
	CellDown   key.Binding
	CellLeft   key.Binding
	CellRight  key.Binding
	NextCell   key.Binding
	NewRow     key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	GotoTop    key.Binding
	GotoBottom key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CellUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CellDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		CellLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		CellRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		NextCell: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next cell"),
		),
		NewRow: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl + n", "new"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),
	}
}

type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
	Error    lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("212")).
			Padding(0, 1),
		Header: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")),
	}
}

func (m *Model) SetStyles(s Styles) {
	m.styles = s

}

// Option is used to set options in New. For example:
//
//	table := New(WithColumns([]Column{{Title: "ID", Width: 10}}))
type Option func(*Model)

func NewModel(opts ...Option) Model {
	m := Model{
		cursorX: 0,
		cursorY: 0,

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(),
		table:  table.New(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.AddRow()

	return m
}

func (m Model) At(row, cell int) string {
	return m.rows[row][cell].input.View()
}

func (m Model) Columns() int {
	return len(m.cols)
}

func (m Model) Rows() int {
	return len(m.rows)
}

func WithColumns(columns []Column) Option {
	return func(m *Model) {
		m.cols = columns
		cols := make([]string, len(m.cols))
		for i := range cols {
			cols[i] = m.cols[i].Title
		}
		m.table.Headers(cols...)
	}
}

func WithRows(rows []Row) Option {
	return func(m *Model) {
		m.rows = rows
	}
}

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

func WithKeyMap(km KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = km
	}
}

func IsOptional(optional bool) Option {
	return func(m *Model) {
		m.optional = optional
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.rows[0][0].input.Focus(), // Focus on the first cell
		textinput.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.CellUp):
			cmd = m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.CellDown):
			cmd = m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.CellLeft):
			cmd = m.MoveLeft(1)
		case key.Matches(msg, m.KeyMap.CellRight):
			cmd = m.MoveRight(1)
		case key.Matches(msg, m.KeyMap.NextCell):
			cmd = m.MoveToNextCell()
		case key.Matches(msg, m.KeyMap.NewRow):
			m.AddRow()
		case key.Matches(msg, m.KeyMap.GotoTop):
			cmd = m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			cmd = m.GotoBottom()
		}
		if cmd != nil {
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	m.rows[m.cursorY][m.cursorX].input, cmd = m.rows[m.cursorY][m.cursorX].input.Update(msg)
	return m, cmd
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
	m.rows[m.cursorY][m.cursorX].input.Blur()
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString(m.table.
		StyleFunc(func(y, x int) lipgloss.Style {
			switch {
			case y == 0:
				return m.styles.Header
			case y == m.cursorY+1 && x == m.cursorX:
				return m.styles.Selected
			default:
				return m.styles.Cell
			}
		}).
		Data(m).
		Render())

	// TODO: Use m.width to manage the max width of the table

	s.WriteRune('\n')
	if errs := m.Errors(); len(errs) != 0 {
		for _, err := range errs {
			s.WriteString(m.styles.Error.Render(fmt.Sprintf("• %s", err.Error())))
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func (m *Model) AddRow() {
	row := make(Row, len(m.cols))
	for i := range row {
		row[i].input = m.newTextInput(m.cols[i])
	}

	m.rows = append(m.rows, row)
}

func (m *Model) RemoveRow(n int) {
	m.rows = append(m.rows[:n], m.rows[n+1:]...)
}

func (m *Model) SetColumns(c []Column) {
	m.cols = c
}

func (m Model) Cursor() (int, int) {
	return m.cursorY, m.cursorX
}

func (m *Model) MoveUp(n int) tea.Cmd {
	return m.Move(-n, 0)
}

func (m *Model) MoveDown(n int) tea.Cmd {
	return m.Move(n, 0)
}

func (m *Model) MoveLeft(n int) tea.Cmd {
	return m.Move(0, -n)
}

func (m *Model) MoveRight(n int) tea.Cmd {
	return m.Move(0, n)
}

func (m *Model) MoveToNextCell() tea.Cmd {
	// If we're not on the last column, move right
	if m.cursorX < len(m.cols)-1 {
		return m.MoveRight(1)
	}

	// else move to the first cell of the next row
	return m.Move(1, -(len(m.cols) - 1))
}

func (m *Model) GotoTop() tea.Cmd {
	return m.Move(-m.cursorY, 0)
}

func (m *Model) GotoBottom() tea.Cmd {
	if m.cursorY == len(m.rows)-1 {
		return nil
	}

	return m.Move(len(m.rows)-1, 0)
}

func (m *Model) Move(y, x int) tea.Cmd {
	if y == 0 && x == 0 {
		return nil
	}

	m.rows[m.cursorY][m.cursorX].input.Blur()
	m.validateCell(m.cursorY, m.cursorX)

	if x != 0 {
		m.cursorX = clamp(m.cursorX+x, 0, len(m.cols)-1)
	}

	if y != 0 {
		if m.cursorY+y >= len(m.rows) {
			for i := 0; i < m.cursorY+y-len(m.rows)+1; i++ {
				m.AddRow()
			}
		}

		if m.cursorY == len(m.rows)-1 && y < 0 && len(m.rows) > 1 {
			isEmpty := true
			for n := 0; n > y; n-- {
				for _, cell := range m.rows[m.cursorY+n] {
					if cell.input.Value() != "" {
						isEmpty = false
						break
					}
				}
				if isEmpty && len(m.rows) > 1 {
					m.RemoveRow(m.cursorY + n)
				}
			}
		}

		m.cursorY = clamp(m.cursorY+y, 0, len(m.rows)-1)
	}

	// Focus on the new cell
	return m.rows[m.cursorY][m.cursorX].input.Focus()
}

func (m Model) Values() [][]string {
	// If the table has only empty cells, return an empty slice
	if m.isEmpty() {
		// TODO: Handle the error here if the table is not optional. Currently it is handled in survey.table component.
		return [][]string{}
	}

	values := make([][]string, len(m.rows))
	for i, row := range m.rows {
		values[i] = make([]string, len(row))
		for j, cell := range row {
			values[i][j] = cell.input.Value()
		}
	}

	return values
}

func (m *Model) Validate() {
	for y := range m.rows {
		for x := range m.rows[y] {
			m.validateCell(y, x)
		}
	}
}

func (m Model) Errors() []error {
	errs := make([]error, 0, len(m.rows)*len(m.cols))
	for y := range m.rows {
		for x := range m.rows[y] {
			if m.rows[y][x].err != nil {
				errs = append(errs, fmt.Errorf("cell (%d, %d): %w", y, x, m.rows[y][x].err))
			}
		}
	}
	return errs
}

func (m *Model) validateCell(y, x int) {
	cell := &m.rows[y][x]
	if m.cols[x].Validators == nil {
		return
	}

	// If the table is optional and the whole table is empty, don't validate the cell.
	// NOTE: Performance-wise it is not smart to check whole table for emptiness on every cell,
	// but it is not a problem since the table is often small.
	if m.optional && m.isEmpty() {
		return
	}

	errs := make([]error, 0, len(m.cols[x].Validators))
	for i := range m.cols[x].Validators {
		err := m.cols[x].Validators[i](cell.input.Value())
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		cell.err = nil
		return
	}

	if len(errs) == 1 {
		cell.err = errs[0]
	}

	errStr := make([]string, len(errs))
	for i := range errs {
		errStr[i] = errs[i].Error()
	}

	cell.err = errors.New(strings.Join(errStr, ", "))
}

// newTextInput initializes a text input which is used inside a cell.
func (m Model) newTextInput(c Column) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""

	ti.Blur()

	return ti
}

func (m Model) isEmpty() bool {
	for y := range m.rows {
		for x := range m.rows[y] {
			if m.rows[y][x].input.Value() != "" {
				return false
			}
		}
	}

	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
