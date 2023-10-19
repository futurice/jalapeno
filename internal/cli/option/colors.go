package option

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/pflag"
)

type Styles struct {
	NoColors bool
	Colors
}

type Colors struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
}

func (opts *Styles) ApplyFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&opts.NoColors, "no-color", false, "If specified, output won't contain any color")
}

func (opts *Styles) Parse() error {
	if opts.NoColors {
		lipgloss.SetColorProfile(termenv.Ascii)
		return nil
	}

	opts.Colors.Primary = lipgloss.Color("#EF4136")
	opts.Colors.Secondary = lipgloss.Color("#26A568")
	return nil
}
