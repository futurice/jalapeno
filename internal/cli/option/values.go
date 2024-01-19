package option

import (
	"errors"

	"github.com/spf13/pflag"
)

type Values struct {
	ReuseOtherSauceValues bool
	NoInput               bool
	CSVDelimiter          rune
	Flags                 []string

	delimiter string
}

func (opts *Values) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&opts.Flags, "set", "s", []string{}, "Set values to be used in the templates. Example: `--set \"MY_VAR=foo\"`")
	fs.StringVar(&opts.delimiter, "delimiter", ",", "Delimiter used when setting table variables")
	fs.BoolVarP(&opts.ReuseOtherSauceValues, "reuse-other-sauce-values", "r", false, "By default each sauce has their own set of values even if the variable names are same in both recipes. Setting this to `true`, values from other sauces will be reused if the variable names match")
	fs.BoolVar(&opts.NoInput, "no-input", false, "If set to true, the program will exit with an error code if it needs to wait for any user input. This is useful when running the program in CI/CD environment")
}

func (opts *Values) Parse() error {
	if opts.delimiter == "" {
		return errors.New("delimiter cannot be empty")
	}
	if len(opts.delimiter) != 1 {
		return errors.New("delimiter can be only one character long")
	}

	opts.CSVDelimiter = rune(opts.delimiter[0])
	return nil
}
