package option

import "github.com/spf13/pflag"

type Output struct {
	OutputPath string
}

func (opts *Output) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.OutputPath, "output", "o", ".", "path where the output files should be created")
}
