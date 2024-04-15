package option

import (
	"time"

	"github.com/spf13/pflag"
)

type Timeout struct {
	Duration time.Duration
}

func (opts *Timeout) ApplyFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&opts.Duration, "timeout", time.Second*10, "Timeout in seconds for the command to run")
}
