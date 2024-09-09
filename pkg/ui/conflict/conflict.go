package conflict

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/diff"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/termenv"
)

const (
	UseOld      int = 1
	UseNew          = 2
	UseDiffFile     = 3
)

type Model struct {
	resolution int
	filePath   string
	fileA      []byte
	fileB      []byte
	diff       diff.Diff
	ready      bool
	err        error
	submitted  bool
	viewport   viewport.Model
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

var _ tea.Model = Model{}

func Solve(in io.Reader, out io.Writer, filePath string, fileA, fileB []byte) ([]byte, error) {
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	p := tea.NewProgram(NewModel(filePath, fileA, fileB), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return []byte{}, err
	} else {
		m, ok := m.(Model)
		if !ok {
			return []byte{}, errors.New("internal error: unexpected model type")
		}

		if m.err != nil {
			return []byte{}, m.err
		}

		return m.Result(), nil
	}
}

func NewModel(filePath string, fileA, fileB []byte) Model {

	fileDiff := diff.New(string(fileA), string(fileB))

	m := Model{
		filePath:   filePath,
		fileA:      fileA,
		fileB:      fileB,
		diff:       fileDiff,
		ready:      false,
		resolution: UseOld,
	}

	m.viewport = viewport.New(20, 20)
	m.viewport.HighPerformanceRendering = false

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) headerView() string {
	title := titleStyle.Render(m.filePath)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, (msg.Height-verticalMarginHeight)/2)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.diff.GetUnifiedDiff())
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = uiutil.ErrUserAborted
			return m, tea.Quit
		case tea.KeyEnter:
			m.submitted = true
			return m, tea.Quit
		case tea.KeyRight:
			m.resolution = min(m.resolution+1, 3)
		case tea.KeyLeft:
			m.resolution = max(m.resolution-1, 1)
		case tea.KeyDown:
			m.viewport.YOffset += 1
			if m.viewport.YOffset >= m.viewport.Height {
				m.viewport.YOffset = m.viewport.Height - 1
			}
		case tea.KeyUp:
			m.viewport.YOffset -= 1
			if m.viewport.YOffset < 0 {
				m.viewport.YOffset = 0
			}
			// case tea.KeyRunes:
			// 	switch string(msg.Runes) {
			// 	case "y", "Y":
			// 		m.answer = true
			// 	case "n", "N":
			// 		m.answer = false
			// 	}
		}
	}
	return m, nil
}

// TODO: Make merge conflict solving more advanced instead of just file override confirmation
func (m Model) View() string {
	var s strings.Builder

	if m.submitted || m.err != nil {
		s.WriteString(fmt.Sprintf("%s: ", m.filePath))
		if m.resolution == UseOld {
			s.WriteString("use old")
		} else if m.resolution == UseNew {
			s.WriteString("use new")
		} else {
			s.WriteString("diff file written")
		}

		s.WriteRune('\n')

		return s.String()
	}

	overwriteOpt := "Overwrite old"
	keepOpt := "Keep old"
	diffOpt := "Write the file with diffs"
	s.WriteString(m.headerView())
	s.WriteString(m.viewport.View())
	s.WriteString(m.footerView())
	s.WriteString("Use up and down arrows to move up and down the diff file.\n Ctrl + down to go the the next conflict and Ctrl + up to go to the previous conflict.")
	s.WriteString(fmt.Sprintf("What to do with file '%s':\n", m.filePath))
	if m.resolution == UseOld {
		s.WriteString(fmt.Sprintf("> %s / %s / %s", lipgloss.NewStyle().Bold(true).Render(keepOpt), overwriteOpt, diffOpt))
	} else if m.resolution == UseNew {
		s.WriteString(fmt.Sprintf("%s / > %s / %s", keepOpt, lipgloss.NewStyle().Bold(true).Render(overwriteOpt), diffOpt))
	} else {
		s.WriteString(fmt.Sprintf("%s / %s / > %s", keepOpt, overwriteOpt, lipgloss.NewStyle().Bold(true).Render(diffOpt)))
	}

	return s.String()
}

func (m Model) Result() []byte {
	if m.resolution == UseOld {
		return m.fileA
	} else if m.resolution == UseNew {
		return m.fileA
	} else {
		return []byte(m.diff.GetConflictResolutionTemplate())
	}
}
