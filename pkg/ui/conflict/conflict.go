package conflict

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	diffFile   []string
	err        error
	submitted  bool
	viewport   viewport.Model
}

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

	m := Model{
		filePath: filePath,
		fileA:    fileA,
		fileB:    fileB,
	}

	m.viewport = viewport.New(20, 20)
	m.viewport.HighPerformanceRendering = false
	m.viewport.SetContent(string(fileA))

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	s.WriteString(fmt.Sprintf("----------------- %s -----------------\n", m.filePath))
	s.WriteString(m.viewport.View())
	s.WriteString("\n")
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
		return []byte("diff file not here!")
	}
}
