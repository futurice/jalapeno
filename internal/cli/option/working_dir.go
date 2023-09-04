package option

import "github.com/spf13/pflag"

type WorkingDirectory struct {
	Dir string
}

func (opts *WorkingDirectory) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Dir, "dir", "d", ".", "Sets working directory")
}
