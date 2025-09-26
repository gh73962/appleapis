package appstoreserver

import (
	"net/http"
	"time"
)

// Client provides a high-level interface to the App Store Server API and JWS verification
type Client struct {
	baseURL        string
	TokenGenerator *TokenGenerator
	httpClient     *http.Client
	userAgent      string
	Verifier       *SignedDataVerifier
}

// New creates a new App Store Server instance using the option pattern
func New(options ...Option) (*Client, error) {
	config := new(ClientConfig)
	for _, option := range options {
		option(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}
	if err := config.Init(); err != nil {
		return nil, err
	}

	c := Client{
		baseURL:   config.Environment.BaseURL(),
		userAgent: "app-store-server-library/go/1.0.0",
	}

	tokenGenerator, err := NewTokenGenerator(config)
	if err != nil {
		return nil, err
	}
	c.TokenGenerator = tokenGenerator

	if config.HTTPClient == nil {
		c.httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	} else {
		c.httpClient = config.HTTPClient
	}

	verifier, err := NewSignedDataVerifier(config)
	if err != nil {
		return nil, err
	}
	c.Verifier = verifier

	return &c, nil
}
