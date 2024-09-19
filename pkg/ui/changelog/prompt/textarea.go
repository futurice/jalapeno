package changelog

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/reflow/wordwrap"
)

type StringModel struct {
	textArea        textarea.Model
	submitted       bool
	showDescription bool
	width           int
	err             error
}

var _ tea.Model = StringModel{}

func NewStringModel() StringModel {
	ti := textarea.New()
	ti.Focus()
	ti.CharLimit = 156

	return StringModel{
		textArea: ti,
		err:      nil,
	}
}

func (m StringModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			err := m.Validate()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.submitted = true
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "?":
				if !m.showDescription {
					m.showDescription = true
					return m, nil
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.textArea.SetWidth(m.width)
	}

	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m StringModel) View() string {
	var s strings.Builder

	if m.submitted {
		if m.textArea.Value() == "" {
			s.WriteString("empty")
		} else {
			s.WriteString(m.textArea.Value())
		}

		return s.String()
	}

	if !m.showDescription {
		s.WriteString(" [type ? for more info]")
	}

	s.WriteRune('\n')
	if m.showDescription {
		s.WriteString(wordwrap.String("Changelog message for version bump\npress Ctrl+S to save", m.width))
		s.WriteRune('\n')
	}

	s.WriteString(m.textArea.View())

	if m.err != nil {
		s.WriteRune('\n')
		errMsg := m.err.Error()
		errMsg = strings.ToUpper(errMsg[:1]) + errMsg[1:]
		s.WriteString(wordwrap.String(errMsg, m.width))
	}

	return s.String()
}

func (m StringModel) Name() string {
	return "Hello world"
}

func (m StringModel) Value() string {
	return m.textArea.Value()
}

func (m StringModel) IsSubmitted() bool {
	return m.submitted
}

func (m StringModel) Validate() error {
	if m.textArea.Value() == "" {
		return util.ErrRequired
	}

	return nil
}
