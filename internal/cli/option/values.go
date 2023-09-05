package option

import "github.com/spf13/pflag"

type Values struct {
	ReuseSauceValues bool
	Flags            []string
}

func (opts *Values) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&opts.Flags, "set", "s", []string{}, "Predefine values to be used in the templates. Example: `--set \"MY_VAR=foo\"`")
	fs.BoolVarP(&opts.ReuseSauceValues, "reuse-sauce-values", "r", false, "By default each sauce has their own set of values even if the variable names are same in both recipes. Setting this to `true` will reuse previous sauce values if the variable name match")
}
