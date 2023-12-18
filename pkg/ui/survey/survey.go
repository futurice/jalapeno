package survey

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/antonmedv/expr"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/survey/prompt"
	"github.com/futurice/jalapeno/pkg/ui/survey/style"
	"github.com/futurice/jalapeno/pkg/ui/util"
	"github.com/muesli/termenv"
)

type SurveyModel struct {
	cursor         int
	submitted      bool
	variables      []recipe.Variable
	existingValues recipe.VariableValues
	prompts        []prompt.Model
	styles         style.Styles
	err            error
}

func NewModel(variables []recipe.Variable, existingValues recipe.VariableValues) SurveyModel {
	model := SurveyModel{
		prompts:        make([]prompt.Model, 0, len(variables)),
		variables:      variables,
		existingValues: existingValues,
		styles:         style.DefaultStyles(),
	}

	p, err := model.createNextPrompt()
	if err != nil {
		model.err = err
	}

	if p != nil {
		model.prompts = append(model.prompts, p)
	}

	return model
}

func (m SurveyModel) Init() tea.Cmd {
	if m.err != nil {
		return tea.Quit
	}

	// Initialize the first prompt (if any)
	if len(m.prompts) > 0 {
		return m.prompts[0].Init()
	}

	return nil
}

func (m SurveyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	// Check if we have already submitted the survey
	if m.submitted {
		return m, nil
	}

	cmds := make([]tea.Cmd, 0, 3)
	submit := func() (tea.Model, tea.Cmd) {
		m.submitted = true
		cmds = append(cmds, tea.Quit)
		return m, tea.Batch(cmds...)
	}

	if len(m.prompts) == 0 {
		return submit()
	}

	lastPrompt := &m.prompts[len(m.prompts)-1]
	promptModel, promptCmd := (*lastPrompt).Update(msg)
	*lastPrompt = promptModel.(prompt.Model)

	if (*lastPrompt).IsSubmitted() {
		cmds = append(cmds, promptCmd)

		// Otherwise, move to the next prompt
		if p, err := m.createNextPrompt(); err != nil {
			m.err = err
			cmds = append(cmds, tea.Quit)
		} else if p == nil {
			return submit()
		} else {
			m.prompts = append(m.prompts, p)
			cmds = append(cmds, p.Init())
		}

		return m, tea.Batch(cmds...)
	}

	return m, promptCmd
}

func (m SurveyModel) View() string {
	var s strings.Builder
	if len(m.prompts) > 0 && !m.submitted && m.err == nil {
		s.WriteString("Provide the following variables:\n\n")
	}

	for i := range m.prompts {
		isLastPrompt := i == len(m.prompts)-1 && len(m.prompts) > 1 && !m.submitted
		if isLastPrompt {
			s.WriteRune('\n')
		}

		s.WriteString(m.prompts[i].View())
		s.WriteRune('\n')
	}

	if m.submitted || m.err != nil {
		s.WriteRune('\n')
	}

	return s.String()
}

func (m SurveyModel) Values() recipe.VariableValues {
	values := make(recipe.VariableValues, len(m.prompts))
	for _, prompt := range m.prompts {
		if prompt.IsSubmitted() {
			values[prompt.Name()] = prompt.Value()
		}
	}

	return values
}

func (m *SurveyModel) createNextPrompt() (prompt.Model, error) {
	if len(m.prompts) > 0 {
		m.cursor++
	}

	if m.cursor >= len(m.variables) {
		return nil, nil
	}

	if p, err := m.createPrompt(m.variables[m.cursor]); err != nil {
		return nil, err
	} else if p == nil {
		return m.createNextPrompt()
	} else {
		return p, nil
	}
}

// createPrompt creates a prompt for the given variable. Returns nil if the variable should be skipped.
func (m SurveyModel) createPrompt(v recipe.Variable) (prompt.Model, error) {
	// Check if variable should be skipped
	if v.If != "" {
		result, err := expr.Eval(v.If, recipeutil.MergeValues(m.existingValues, m.Values()))
		if err != nil {
			return nil, fmt.Errorf("error when evaluating variable \"%s\" 'if' expression: %w", v.Name, err)
		}
		variableShouldBePrompted, ok := result.(bool)
		if !ok {
			return nil, fmt.Errorf("result of 'if' expression of variable \"%s\" was not a boolean value, was %T instead", v.Name, result)
		}

		if !variableShouldBePrompted {
			return nil, nil
		}
	}

	var p prompt.Model
	switch {
	case len(v.Options) != 0:
		p = prompt.NewSelectModel(v, m.styles)
	case v.Confirm:
		p = prompt.NewConfirmModel(v, m.styles)
	case len(v.Columns) > 0:
		p = prompt.NewTableModel(v, m.styles)
	default:
		p = prompt.NewStringModel(v, m.styles)
	}

	return p, nil
}

// PromptUserForValues prompts the user for values for the given variables
func PromptUserForValues(in io.Reader, out io.Writer, variables []recipe.Variable, existingValues recipe.VariableValues) (recipe.VariableValues, error) {
	// https://github.com/charmbracelet/lipgloss/issues/73#issuecomment-1144921037
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	p := tea.NewProgram(NewModel(variables, existingValues), tea.WithInput(in), tea.WithOutput(out))
	if m, err := p.Run(); err != nil {
		return nil, err
	} else {
		survey, ok := m.(SurveyModel)
		if !ok {
			return nil, errors.New("internal error: unexpected model type")
		}
		if survey.err != nil {
			return nil, survey.err
		}

		if survey.submitted {
			return m.(SurveyModel).Values(), nil
		}

		return nil, util.ErrUserAborted
	}
}
