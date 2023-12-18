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
	Value     bool
	filePath  string
	err       error
	submitted bool
}

var _ tea.Model = Model{}

func Solve(in io.Reader, out io.Writer, filePath string) (bool, error) {
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	p := tea.NewProgram(NewModel(filePath), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return false, err
	} else {
		m, ok := m.(Model)
		if !ok {
			return false, errors.New("internal error: unexpected model type")
		}

		if m.err != nil {
			return false, m.err
		}

		return m.Value, nil
	}
}

func NewModel(filePath string) Model {
	return Model{
		filePath: filePath,
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
			m.Value = true
		case tea.KeyLeft:
			m.Value = false
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "y", "Y":
				m.Value = true
			case "n", "N":
				m.Value = false
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.submitted || m.err != nil {
		return ""
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("Override file '%s':\n", m.filePath))
	if m.Value {
		s.WriteString(fmt.Sprintf("> No/%s", lipgloss.NewStyle().Bold(true).Render("Yes")))
	} else {
		s.WriteString(fmt.Sprintf("> %s/Yes", lipgloss.NewStyle().Bold(true).Render("No")))
	}

	return s.String()
}
