package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
	"github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/reflow/wordwrap"
)

type StringModel struct {
	variable        recipe.Variable
	textInput       textinput.Model
	styles          style.Styles
	submitted       bool
	showDescription bool
	width           int
	err             error
}

var _ Model = StringModel{}

func NewStringModel(v recipe.Variable, styles style.Styles) StringModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	if v.Default != "" {
		ti.SetValue(v.Default)
	}

	return StringModel{
		variable:  v,
		textInput: ti,
		err:       nil,
		styles:    styles,
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
		case tea.KeyEnter:
			if err := m.Validate(); err != nil {
				m.err = err
				return m, nil
			}
			m.submitted = true
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
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StringModel) View() string {
	var s strings.Builder
	s.WriteString(m.styles.VariableName.Render(m.variable.Name))

	if m.submitted {
		s.WriteString(": ")

		if m.textInput.Value() == "" {
			s.WriteString(m.styles.HelpText.Render("empty"))
		} else {
			s.WriteString(m.textInput.Value())
		}

		return s.String()
	}

	if m.variable.Optional {
		s.WriteString(m.styles.HelpText.Render(" (optional)"))
	}

	if m.variable.Description != "" && !m.showDescription {
		s.WriteString(m.styles.HelpText.Render(" [type ? for more info]"))
	}

	s.WriteRune('\n')
	if m.showDescription {
		s.WriteString(wordwrap.String(m.variable.Description, m.width))
		s.WriteRune('\n')
	}

	s.WriteString(m.textInput.View())

	if m.err != nil {
		s.WriteRune('\n')
		errMsg := m.err.Error()
		errMsg = strings.ToUpper(errMsg[:1]) + errMsg[1:]
		s.WriteString(wordwrap.String(m.styles.ErrorText.Render(errMsg), m.width))
	}

	return s.String()
}

func (m StringModel) Name() string {
	return m.variable.Name
}

func (m StringModel) Value() interface{} {
	return m.textInput.Value()
}

func (m StringModel) IsSubmitted() bool {
	return m.submitted
}

func (m StringModel) Validate() error {
	if !m.variable.Optional && m.textInput.Value() == "" {
		return util.ErrRequired
	}

	for _, v := range m.variable.Validators {
		if v.Pattern != "" {
			validatorFunc, err := v.CreateValidatorFunc()
			if err != nil {
				return fmt.Errorf("validator function create failed: %s", err)
			}
			if err := validatorFunc(m.textInput.Value()); err != nil {
				return fmt.Errorf("%w: %s", util.ErrRegExFailed, err)
			}
		}
	}

	return nil
}
