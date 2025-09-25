package appstoreserver

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// client represents the App Store Server API client
type client struct {
	baseURL        string
	tokenGenerator *TokenGenerator
	httpClient     *http.Client
	userAgent      string
	verifier       *SignedDataVerifier
}

// newClient creates a new App Store Server API client
func newClient(environment Environment, tokenGenerator *TokenGenerator, httpClient *http.Client) (*client, error) {
	var baseURL string
	switch environment {
	case EnvironmentProduction:
		baseURL = ProductionBaseURL
	case EnvironmentSandbox:
		baseURL = SandboxBaseURL
	case EnvironmentLocalTesting:
		baseURL = LocalTestingBaseURL
	default:
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &client{
		baseURL:        baseURL,
		tokenGenerator: tokenGenerator,
		httpClient:     httpClient,
		userAgent:      "app-store-server-library/go/1.0.0",
	}, nil
}

// Client provides a high-level interface to the App Store Server API and JWS verification
type Client struct {
	*client
}

// New creates a new App Store Server instance using the option pattern
func New(options ...ClientOption) (*Client, error) {
	config := new(ClientConfig)

	for _, option := range options {
		option(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	privateKey, err := ParsePrivateKeyFromPEM(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	tokenGenerator := NewTokenGenerator(privateKey, config.KeyID, config.IssuerID, config.BundleID)

	client, err := newClient(config.Environment, tokenGenerator, config.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	var rootCerts [][]byte
	// If no root certificates provided, use Apple's default root certificates
	if len(config.RootCertificates) == 0 {
		for _, v := range []string{AppleRootCAURL, AppleRootCAG2URL, AppleRootCAG3URL} {
			resp, err := http.Get(v)
			if err != nil {
				return nil, fmt.Errorf("failed to download %s: %w", v, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to download %s: HTTP %d", v, resp.StatusCode)
			}

			for _, v := range []string{AppleRootCAURL, AppleRootCAG2URL, AppleRootCAG3URL} {
				resp, err := http.Get(v)
				if err != nil {
					return nil, fmt.Errorf("failed to download %s: %w", v, err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("failed to download %s: HTTP %d", v, resp.StatusCode)
				}

				// 修改：使用io.ReadAll替代手动缓冲读取
				certData, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read certificate data from %s: %w", v, err)
				}

				rootCerts = append(rootCerts, certData)
			}

		}

	} else {
		rootCerts = config.RootCertificates
	}

	verifier, err := NewSignedDataVerifier(
		rootCerts,
		config.EnableOnlineChecks,
		config.Environment,
		config.BundleID,
		config.AppAppleID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create verifier: %w", err)
	}

	client.verifier = verifier
	return &Client{
		client: client,
	}, nil
}
