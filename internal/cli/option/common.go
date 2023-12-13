package option

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/pflag"
)

type Common struct {
	Debug    bool
	Verbose  bool
	NoColors bool
	Colors
}

type Colors struct {
	Green lipgloss.Style
	Red   lipgloss.Style
}

func (opts *Common) ApplyFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&opts.Debug, "debug", false, "Debug mode")
	fs.BoolVarP(&opts.Verbose, "verbose", "v", false, "Verbose output")
	fs.BoolVar(&opts.NoColors, "no-color", false, "If specified, output won't contain any color")
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
