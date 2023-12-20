package conflict

import (
	"errors"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/termenv"
)

type Model struct {
	answer    bool
	filePath  string
	fileA     []byte
	fileB     []byte
	err       error
	submitted bool
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
	return Model{
		filePath: filePath,
		fileA:    fileA,
		fileB:    fileB,
	}
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
			m.answer = true
		case tea.KeyLeft:
			m.answer = false
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "y", "Y":
				m.answer = true
			case "n", "N":
				m.answer = false
			}
		}
	}
	return m, nil
}

// TODO: Make merge conflict solving more advanced instead of just file override confirmation
func (m Model) View() string {
	var s strings.Builder
	if m.submitted || m.err != nil {
		s.WriteString(fmt.Sprintf("%s: ", m.filePath))
		if m.answer {
			s.WriteString("override")
		} else {
			s.WriteString("keep")
		}

		return s.String()
	}

	s.WriteString(fmt.Sprintf("Override file '%s':\n", m.filePath))
	if m.answer {
		s.WriteString(fmt.Sprintf("> No/%s", lipgloss.NewStyle().Bold(true).Render("Yes")))
	} else {
		s.WriteString(fmt.Sprintf("> %s/Yes", lipgloss.NewStyle().Bold(true).Render("No")))
	}

	return s.String()
}

func (m Model) Result() []byte {
	if m.answer {
		return m.fileB
	} else {
		return m.fileA
	}
}
