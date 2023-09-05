package oci

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/futurice/jalapeno/pkg/oci/credential"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type Repository struct {
	Reference   string
	PlainHTTP   bool
	Credentials Credentials
	TLS         TLSConfig
}

type Credentials struct {
	Username      string
	Password      string
	DockerConfigs []string
}

type TLSConfig struct {
	CACertFilePath string
	Insecure       bool
}

// NewRepository assembles an ORAS remote repository.
func NewRepository(opts Repository) (*remote.Repository, error) {
	repo, err := remote.NewRepository(opts.Reference)
	if err != nil {
		return nil, err
	}

	repo.PlainHTTP = opts.PlainHTTP

	tlsConfig := &tls.Config{
		InsecureSkipVerify: opts.TLS.Insecure,
	}

	if opts.TLS.CACertFilePath != "" {
		tlsConfig.RootCAs, err = loadCertPool(opts.TLS.CACertFilePath)
		if err != nil {
			return nil, err
		}
	}

	client := &auth.Client{
		Client: &http.Client{
			// default value are derived from http.DefaultTransport
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       tlsConfig,
			},
		},
		Header: http.Header{
			"User-Agent": []string{fmt.Sprintf("jalapeno/%s", "0.0.1")}, // TODO: Get real version number
		},
		Cache: auth.NewCache(),
	}

	if client.Credential, err = GetCredentials(repo.Reference.Registry, opts.Credentials); err != nil {
		return nil, err
	}

	repo.Client = client

	return repo, err
}

func GetCredentials(registry string, creds Credentials) (func(context.Context, string) (auth.Credential, error), error) {
	cred := credential.Credential(creds.Username, creds.Password)
	if cred != auth.EmptyCredential {
		return func(ctx context.Context, s string) (auth.Credential, error) {
			return cred, nil
		}, nil
	} else {
		store, err := credential.NewStore(creds.DockerConfigs...)
		if err != nil {
			return nil, err
		}

		// For a user case with a registry from 'docker.io', the hostname is "registry-1.docker.io"
		// According to the the behavior of Docker CLI,
		// credential under key "https://index.docker.io/v1/" should be provided
		if registry == "docker.io" {
			return func(ctx context.Context, hostname string) (auth.Credential, error) {
				if hostname == "registry-1.docker.io" {
					hostname = "https://index.docker.io/v1/"
				}
				return store.Credential(ctx, hostname)
			}, nil
		} else {
			return store.Credential, nil
		}
	}
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
