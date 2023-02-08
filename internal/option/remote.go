package option

import "github.com/spf13/pflag"

type Remote struct {
	CACertFilePath    string
	PlainHTTP         bool
	Insecure          bool
	Configs           []string
	Username          string
	PasswordFromStdin bool
	Password          string
}

func (opts *Remote) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Username, "username", "u", "", "registry username")
	fs.StringVarP(&opts.Password, "password", "p", "", "registry password or identity token")
	fs.BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	fs.BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "allow insecure connections to registry without SSL check")
	fs.StringVarP(&opts.CACertFilePath, "ca-file", "", "", "server certificate authority file for the remote registry")
}
