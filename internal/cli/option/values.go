package option

import (
	"errors"

	"github.com/spf13/pflag"
)

type Values struct {
	ReuseSauceValues bool
	CSVDelimiter     rune
	Flags            []string

	delimiter string
}

func (opts *Values) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&opts.Flags, "set", "s", []string{}, "Predefine values to be used in the templates. Example: `--set \"MY_VAR=foo\"`")
	fs.StringVar(&opts.delimiter, "delimiter", ",", "Delimiter used when setting table variables")
	fs.BoolVarP(&opts.ReuseSauceValues, "reuse-sauce-values", "r", false, "By default each sauce has their own set of values even if the variable names are same in both recipes. Setting this to `true` will reuse previous sauce values if the variable name match")
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
