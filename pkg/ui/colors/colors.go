package colors

import "github.com/charmbracelet/lipgloss"

var (
	Red       = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4136"))
	Green     = lipgloss.NewStyle().Foreground(lipgloss.Color("#26A568"))
	Highlight = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#00EEFF", Dark: "#FFFF00"})
)
