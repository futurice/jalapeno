package option

import "github.com/spf13/pflag"

type Values struct {
	Flags []string
}

func (opts *Values) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&opts.Flags, "set", "s", []string{}, "Predefine values to be used in the templates. Example: `--set \"MY_VAR=foo\"`")
}
