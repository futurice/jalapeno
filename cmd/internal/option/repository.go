package option

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/futurice/jalapeno/cmd/internal/credential"
	"github.com/spf13/pflag"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type Repository struct {
	CACertFilePath string
	PlainHTTP      bool
	Insecure       bool
	Configs        []string
}

func (opts *Repository) ApplyFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	fs.BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "allow insecure connections to registry without SSL check")
	fs.StringVarP(&opts.CACertFilePath, "ca-file", "", "", "server certificate authority file for the remote registry")
}

// NewRegistry assembles a oras remote registry.
func (opts *Repository) NewRegistry(hostname string, common Common) (reg *remote.Registry, err error) {
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

// NewRepository assembles a oras remote repository.
func (opts *Repository) NewRepository(reference string, common Common) (repo *remote.Repository, err error) {
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

func (opts *Repository) tlsConfig() (*tls.Config, error) {
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
func (opts *Repository) authClient(registry string, debug bool) (client *auth.Client, err error) {
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

	return
}
