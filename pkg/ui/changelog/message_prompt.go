package changelog

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/muesli/reflow/wordwrap"
)

type MessageModel struct {
	textArea textarea.Model
	width    int
	err      error
}

var _ tea.Model = MessageModel{}

func NewStringModel() MessageModel {
	ti := textarea.New()
	ti.Focus()
	ti.SetHeight(5)
	ti.CharLimit = 156

	return MessageModel{
		textArea: ti,
		err:      nil,
	}
}

func (m MessageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MessageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlS:
			m.err = m.Validate()
			if m.err == nil {
				return m, tea.Quit
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.textArea.SetWidth(m.width)
	}

	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m MessageModel) View() string {
	var s strings.Builder

	s.WriteString(wordwrap.String("Write the changelog message for the new version", m.width))
	s.WriteString("\nPress Ctrl+S to save")
	s.WriteRune('\n')

	s.WriteString(m.textArea.View())

	if m.err != nil {
		s.WriteString("\n\n")
		s.WriteString(colors.Red.Render(fmt.Sprintf("Error: %s", m.err.Error())))
	}

	s.WriteString("\n\n")

	return s.String()
}

func (m MessageModel) Value() string {
	return strings.TrimSpace(m.textArea.Value())
}

func (m MessageModel) Validate() error {
	if m.textArea.Value() == "" {
		return errors.New("changelog message cannot be empty")
	}

	return nil
}
