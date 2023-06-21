package option

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/futurice/jalapeno/internal/cli/internal/credential"
	"github.com/spf13/pflag"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
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
	fs.StringVarP(&opts.Username, "username", "u", "", "registry username")
	fs.StringVarP(&opts.Password, "password", "p", "", "registry password or identity token")
	fs.BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	fs.BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "allow insecure connections to registry without SSL check")
	fs.StringVarP(&opts.CACertFilePath, "ca-file", "", "", "server certificate authority file for the remote registry")
	fs.StringArrayVarP(&opts.Configs, "registry-config", "", nil, "`path` of the authentication file")
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

// NewRegistry assembles an ORAS remote registry.
func (opts *OCIRepository) NewRegistry(hostname string, common Common) (reg *remote.Registry, err error) {
	reg, err = remote.NewRegistry(hostname)
	if err != nil {
		return nil, err
	}
	hostname = reg.Reference.Registry
	reg.PlainHTTP = opts.PlainHTTP
	if reg.Client, err = opts.authClient(hostname, common.Debug); err != nil {
		return nil, err
	}
	return
}

// NewRepository assembles an ORAS remote repository.
func (opts *OCIRepository) NewRepository(reference string, common Common) (repo *remote.Repository, err error) {
	repo, err = remote.NewRepository(reference)
	if err != nil {
		return nil, err
	}
	hostname := repo.Reference.Registry
	repo.PlainHTTP = opts.PlainHTTP
	if repo.Client, err = opts.authClient(hostname, common.Debug); err != nil {
		return nil, err
	}
	return
}

func (opts *OCIRepository) tlsConfig() (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: opts.Insecure,
	}
	if opts.CACertFilePath != "" {
		var err error
		config.RootCAs, err = loadCertPool(opts.CACertFilePath)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func loadCertPool(path string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if ok := pool.AppendCertsFromPEM(pemBytes); !ok {
		return nil, errors.New("Failed to load certificate in file: " + path)
	}
	return pool, nil
}

// authClient assembles a oras auth client.
func (opts *OCIRepository) authClient(registry string, debug bool) (client *auth.Client, err error) {
	config, err := opts.tlsConfig()
	if err != nil {
		return nil, err
	}

	client = &auth.Client{
		Client: &http.Client{
			// default value are derived from http.DefaultTransport
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       config,
			},
		},
		Cache: auth.NewCache(),
	}
	client.SetUserAgent("jalapeno/0.0.1") // TODO: Get real version number

	cred := opts.Credential()
	if cred != auth.EmptyCredential {
		client.Credential = func(ctx context.Context, s string) (auth.Credential, error) {
			return cred, nil
		}
	} else {
		store, err := credential.NewStore(opts.Configs...)
		if err != nil {
			return nil, err
		}
		// For a user case with a registry from 'docker.io', the hostname is "registry-1.docker.io"
		// According to the the behavior of Docker CLI,
		// credential under key "https://index.docker.io/v1/" should be provided
		if registry == "docker.io" {
			client.Credential = func(ctx context.Context, hostname string) (auth.Credential, error) {
				if hostname == "registry-1.docker.io" {
					hostname = "https://index.docker.io/v1/"
				}
				return store.Credential(ctx, hostname)
			}
		} else {
			client.Credential = store.Credential
		}
	}

	return
}

// Credential returns a credential based on the remote options.
func (opts *OCIRepository) Credential() auth.Credential {
	return credential.Credential(opts.Username, opts.Password)
}
