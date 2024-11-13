package changelog

import (
	"errors"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/futurice/jalapeno/pkg/ui/util"
)

const (
	Patch = "patch"
	Minor = "minor"
	Major = "major"
)

type Changelog struct {
	Increment string
	Msg       string
}

func RunChangelog(in io.Reader, out io.Writer) (Changelog, error) {
	verInc, err := runVersionPrompt(in, out)

	if err != nil {
		return Changelog{}, fmt.Errorf("failed to get version type: %w", err)
	}

	logmsg, err := runMessagePrompt(in, out)

	if err != nil {
		return Changelog{}, fmt.Errorf("failed to get log message: %w", err)
	}

	changelog := Changelog{
		Increment: verInc,
		Msg:       logmsg,
	}

	return changelog, nil
}

func runVersionPrompt(in io.Reader, out io.Writer) (string, error) {
	options := []string{Patch, Minor, Major}

	p := tea.NewProgram(NewSelectModel(options), tea.WithInput(in), tea.WithOutput(out))

	if m, err := p.Run(); err != nil {
		return "", err
	} else {
		sel, ok := m.(VersionModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}
		value := sel.Value()
		if value == "" {
			return "", util.ErrUserAborted
		}

		return value, nil
	}
}

func runMessagePrompt(in io.Reader, out io.Writer) (string, error) {
	p := tea.NewProgram(NewStringModel(), tea.WithInput(in), tea.WithOutput(out))

	if m, err := p.Run(); err != nil {
		return "", err
	} else {
		txt, ok := m.(MessageModel)
		if !ok {
			return "", errors.New("internal error: unexpected model type")
		}

		value := txt.Value()
		if value == "" {
			return "", util.ErrUserAborted
		}

		return txt.Value(), nil
	}
}
