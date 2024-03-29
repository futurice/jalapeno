package style

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	VariableName lipgloss.Style
	ErrorText    lipgloss.Style
	HelpText     lipgloss.Style
	Bold         lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		VariableName: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")),
		ErrorText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")),
		HelpText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999")),
		Bold: lipgloss.NewStyle().Bold(true),
	}
}
