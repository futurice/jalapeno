package option

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/pflag"
)

type Common struct {
	NoColors bool
	NoInput  bool
}

func (opts *Common) ApplyFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&opts.NoColors, "no-color", false, "If specified, output won't contain any color")
	fs.BoolVar(&opts.NoInput, "no-input", false, "If set to true, the program will exit with an error code if it needs to wait for any user input. This is useful when running the program in CI/CD environment")
}

func (opts *Common) Parse() error {
	if opts.NoColors {
		lipgloss.SetColorProfile(termenv.Ascii)
		return nil
	}

	return nil
}
