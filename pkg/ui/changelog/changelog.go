package changelog

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	changelog "github.com/futurice/jalapeno/pkg/ui/changelog/prompt"
)

func RunChangelog() ([]string, error) {
	var changelog []string
	verInc, err := RunSelectPrompt()

	if err != nil {
		return []string{}, errors.New("failed to get version type")
	}

	logmsg, err := RunTextAreaPrompt()

	if err != nil {
		return []string{}, errors.New("failed to get log message")
	}

	changelog = append(changelog, verInc)
	changelog = append(changelog, logmsg)

	return changelog, nil
}

func RunSelectPrompt() (string, error) {
	options := []string{"patch", "minor", "major"}

	p := tea.NewProgram(changelog.NewSelectModel(options))

	if m, err := p.Run(); err != nil {
		log.Fatal(err)
	} else {
		sel, ok := m.(changelog.SelectModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}
		return sel.Value(), nil
	}
	return "", nil
}

func RunTextAreaPrompt() (string, error) {
	p := tea.NewProgram(changelog.NewStringModel())

	if m, err := p.Run(); err != nil {
		log.Fatal(err)
	} else {
		txt, ok := m.(changelog.StringModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}
		return txt.Value(), nil
	}
	return "", nil
}
