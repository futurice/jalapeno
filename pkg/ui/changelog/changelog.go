package changelog

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	changelog "github.com/futurice/jalapeno/pkg/ui/changelog/prompt"
)

type Changelog struct {
	Increment string
	Msg       string
}

func RunChangelog() (Changelog, error) {
	verInc, err := runSelectPrompt()

	if err != nil {
		return Changelog{}, errors.New("failed to get version type")
	}

	logmsg, err := runTextAreaPrompt()

	if err != nil {
		return Changelog{}, errors.New("failed to get log message")
	}

	changelog := Changelog{
		Increment: verInc,
		Msg:       logmsg,
	}

	return changelog, nil
}

func runSelectPrompt() (string, error) {
	options := []string{"patch", "minor", "major"}

	p := tea.NewProgram(changelog.NewSelectModel(options))

	if m, err := p.Run(); err != nil {
		return "", err
	} else {
		sel, ok := m.(changelog.SelectModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}
		return sel.Value(), nil
	}
}

func runTextAreaPrompt() (string, error) {
	p := tea.NewProgram(changelog.NewStringModel())

	if m, err := p.Run(); err != nil {
		return "", err
	} else {
		txt, ok := m.(changelog.StringModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}
		return txt.Value(), nil
	}
}
