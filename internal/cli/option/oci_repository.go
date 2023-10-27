package option

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/spf13/pflag"
)

type OCIRepository struct {
	CACertFilePath    string
	PlainHTTP         bool
	Insecure          bool
	Configs           []string
	Username          string
	PasswordFromStdin bool
	Password          string
}

func (opts *OCIRepository) ApplyFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&opts.Username, "username", "u", "", "Registry username")
	fs.StringVarP(&opts.Password, "password", "p", "", "Registry password or identity token")
	fs.BoolVarP(&opts.Insecure, "insecure", "", false, "Allow connections to SSL registry without certs")
	fs.BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "Allow insecure connections to registry without SSL check")
	fs.StringVarP(&opts.CACertFilePath, "ca-file", "", "", "Server certificate authority file for the remote registry")
	fs.StringArrayVarP(&opts.Configs, "registry-config", "", nil, "Path of the authentication file")
}

// Parse tries to read password with optional cmd prompt.
func (opts *OCIRepository) Parse() error {
	return opts.readPassword()
}

// readPassword tries to read password with optional cmd prompt.
func (opts *OCIRepository) readPassword() (err error) {
	if opts.Password != "" {
		fmt.Fprintln(os.Stderr, "WARNING! Using --password via the CLI is insecure. Use --password-stdin.")
	} else if opts.PasswordFromStdin {
		// Prompt for credential
		password, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		opts.Password = strings.TrimSuffix(string(password), "\n")
		opts.Password = strings.TrimSuffix(opts.Password, "\r")
	}
	return nil
}

func (opts *OCIRepository) Repository(url string) oci.Repository {
	return oci.Repository{
		Reference: strings.TrimPrefix(url, "oci://"),
		PlainHTTP: opts.PlainHTTP,
		Credentials: oci.Credentials{
			Username:      opts.Username,
			Password:      opts.Password,
			DockerConfigs: opts.Configs,
		},
		TLS: oci.TLSConfig{
			CACertFilePath: opts.CACertFilePath,
			Insecure:       opts.Insecure,
		},
	}
}
