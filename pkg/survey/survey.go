package survey

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type SurveyModel struct {
	cursor    int
	variables []recipe.Variable
	prompts   []PromptModel
}

type PromptModel interface {
	tea.Model
	Value() interface{}
}

var _ tea.Model = PromptModel(nil)

type FocusMsg struct{}
type BlurMsg struct{}

func NewSurveyModel(variables []recipe.Variable) SurveyModel {
	model := SurveyModel{
		prompts:   make([]PromptModel, 0, len(variables)),
		variables: variables,
	}

	for _, variable := range variables {
		var prompt PromptModel
		switch {
		case len(variable.Options) != 0:
			// prompt = NewSelectModel() // TODO
			prompt = NewStringPromptModel(variable)
		case variable.Confirm:
			// prompt = NewConfirmModel() // TODO
			prompt = NewStringPromptModel(variable)
		case len(variable.Columns) > 0:
			// prompt = NewTableModel() // TODO
			prompt = NewStringPromptModel(variable)
		default:
			prompt = NewStringPromptModel(variable)
		}
		model.prompts = append(model.prompts, prompt)
	}

	return model
}

func (m SurveyModel) Init() tea.Cmd {
	// Initialize the first prompt
	return m.prompts[0].Init()
}

func (m SurveyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: if property
	// TODO: regex validate property

	var updatedModel tea.Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds := make([]tea.Cmd, 2)
			// Unfocus the current prompt
			updatedModel, cmds[0] = m.prompts[m.cursor].Update(BlurMsg{})
			m.prompts[m.cursor] = updatedModel.(PromptModel)

			// Check if we're on the last prompt
			if m.cursor == len(m.prompts)-1 {
				cmds[1] = tea.Quit
				return m, tea.Batch(cmds...)
			}

			// Otherwise, move to the next prompt
			m.cursor++
			cmds[1] = m.prompts[m.cursor].Init()
			return m, tea.Batch(cmds...)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var promptCmd tea.Cmd
	updatedModel, promptCmd = m.prompts[m.cursor].Update(msg)
	m.prompts[m.cursor] = updatedModel.(PromptModel)
	return m, promptCmd
}

func (m SurveyModel) View() (s string) {
	s += "Provide the following variables:\n\n"

	for i := 0; i <= m.cursor; i++ {
		s += m.prompts[i].View()
		s += "\n\n"
	}

	return
}

func (m SurveyModel) Values() recipe.VariableValues {
	values := make(recipe.VariableValues, len(m.prompts))
	for i, prompt := range m.prompts {
		switch prompt := prompt.(type) {
		case PromptModel:
			values[m.variables[i].Name] = prompt.Value()
		}
	}

	return values
}

// PromptUserForValues prompts the user for values for the given variables
func PromptUserForValues(in io.Reader, out io.Writer, variables []recipe.Variable, existingValues recipe.VariableValues) (recipe.VariableValues, error) {
	p := tea.NewProgram(NewSurveyModel(variables), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return nil, err
	} else {
		return m.(SurveyModel).Values(), nil
	}
}
