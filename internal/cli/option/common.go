package option

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/pflag"
)

type Common struct {
	Debug    bool
	NoColors bool
	NoInput  bool
	Colors
}

type Colors struct {
	Green lipgloss.Style
	Red   lipgloss.Style
}

func (opts *Common) ApplyFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&opts.Debug, "debug", false, "Debug mode")
	fs.BoolVar(&opts.NoColors, "no-color", false, "If specified, output won't contain any color")
	fs.BoolVar(&opts.NoInput, "no-input", false, "If set to true, the program will exit with an error code if it needs to wait for any user input. This is useful when running the program in CI/CD environment")
}

func (opts *Common) Parse() error {
	if opts.NoColors {
		lipgloss.SetColorProfile(termenv.Ascii)
		return nil
	}

	opts.Colors = Colors{
		Red:   lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4136")),
		Green: lipgloss.NewStyle().Foreground(lipgloss.Color("#26A568")),
	}

	return nil
}
