package appstoreserver

import (
	"fmt"
	"net/http"
	"time"
)

// client represents the App Store Server API client
type client struct {
	baseURL        string
	tokenGenerator *TokenGenerator
	httpClient     *http.Client
	userAgent      string
}

// newClient creates a new App Store Server API client
func newClient(environment Environment, tokenGenerator *TokenGenerator) (*client, error) {
	var baseURL string
	switch environment {
	case EnvironmentProduction:
		baseURL = ProductionBaseURL
	case EnvironmentSandbox:
		baseURL = SandboxBaseURL
	default:
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}

	return &client{
		baseURL:        baseURL,
		tokenGenerator: tokenGenerator,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "app-store-server-library/go/1.0.0",
	}, nil
}

// Client provides a high-level interface to the App Store Server API and JWS verification
type Client struct {
	client   *client
	verifier *SignedDataVerifier
}

// New creates a new App Store Server instance using the option pattern
func New(options ...ClientOption) (*Client, error) {
	config := new(ClientConfig)

	for _, option := range options {
		option(config)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}
	// Parse private key
	privateKey, err := ParsePrivateKeyFromPEM(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create token generator
	tokenGenerator := NewTokenGenerator(privateKey, config.KeyID, config.IssuerID, config.BundleID)

	// Create API client
	client, err := newClient(config.Environment, tokenGenerator)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Create JWS verifier
	var rootCerts [][]byte = config.RootCertificates
	// TODO: If no root certificates provided, use Apple's default root certificates

	verifierOptions := []SignedDataVerifierOption{
		WithRootCertificates(rootCerts),
		WithEnvironment(config.Environment),
		WithBundleID(config.BundleID),
	}

	// Add AppAppleID only if it's provided (not nil)
	if config.AppAppleID != nil {
		verifierOptions = append(verifierOptions, WithAppAppleID(*config.AppAppleID))
	}

	verifier, err := NewSignedDataVerifier(verifierOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWS verifier: %w", err)
	}

	return &Client{
		client:   client,
		verifier: verifier,
	}, nil
}
