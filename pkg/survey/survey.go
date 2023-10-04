package survey

import (
	"errors"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/survey/prompt"
	"github.com/futurice/jalapeno/pkg/survey/util"
)

type SurveyModel struct {
	cursor    int
	submitted bool
	variables []recipe.Variable
	prompts   []prompt.Model
}

var (
	ErrUserAborted = errors.New("user aborted")
)

func NewSurveyModel(variables []recipe.Variable) SurveyModel {
	model := SurveyModel{
		prompts:   make([]prompt.Model, 0, len(variables)),
		variables: variables,
	}

	for _, variable := range variables {
		var p prompt.Model
		switch {
		case len(variable.Options) != 0:
			// prompt = NewSelectModel() // TODO
			p = prompt.NewStringModel(variable)
		case variable.Confirm:
			// prompt = NewConfirmModel() // TODO
			p = prompt.NewStringModel(variable)
		case len(variable.Columns) > 0:
			// prompt = NewTableModel() // TODO
			p = prompt.NewStringModel(variable)
		default:
			p = prompt.NewStringModel(variable)
		}
		model.prompts = append(model.prompts, p)
	}

	return model
}

func (m SurveyModel) Init() tea.Cmd {
	// Initialize the first prompt
	return m.prompts[0].Init()
}

func (m SurveyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: if property

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	promptModel, promptCmd := m.prompts[m.cursor].Update(msg)
	m.prompts[m.cursor] = promptModel.(prompt.Model)

	if m.prompts[m.cursor].IsSubmitted() {
		cmds := make([]tea.Cmd, 0, 3)
		cmds = append(cmds, promptCmd)

		// Unfocus the current prompt
		promptModel, promptCmd = m.prompts[m.cursor].Update(util.Blur())
		m.prompts[m.cursor] = promptModel.(prompt.Model)
		cmds = append(cmds, promptCmd)

		// Check if we're on the last prompt
		if m.cursor == len(m.prompts)-1 {
			m.submitted = true
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		// Otherwise, move to the next prompt
		m.cursor++
		cmds = append(cmds, m.prompts[m.cursor].Init())
		return m, tea.Batch(cmds...)
	}

	return m, promptCmd
}

func (m SurveyModel) View() (s string) {
	s += "Provide the following variables:\n\n"
	for i := 0; i <= m.cursor; i++ {
		cursorIsInLastVisiblePrompt := i == m.cursor && i != 0 && !m.submitted
		if cursorIsInLastVisiblePrompt {
			s += "\n"
		}

		s += m.prompts[i].View()

		if !cursorIsInLastVisiblePrompt {
			s += "\n"
		}
	}

	return
}

func (m SurveyModel) Values() recipe.VariableValues {
	values := make(recipe.VariableValues, len(m.prompts))
	for i, prompt := range m.prompts {
		values[m.variables[i].Name] = prompt.Value()
	}

	return values
}

// PromptUserForValues prompts the user for values for the given variables
func PromptUserForValues(in io.Reader, out io.Writer, variables []recipe.Variable, existingValues recipe.VariableValues) (recipe.VariableValues, error) {
	p := tea.NewProgram(NewSurveyModel(variables), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return nil, err
	} else {
		survey := m.(SurveyModel)
		if survey.submitted {
			return m.(SurveyModel).Values(), nil
		}

		return nil, ErrUserAborted
	}
}
