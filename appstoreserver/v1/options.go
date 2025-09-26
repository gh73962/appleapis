package appstoreserver

import (
	"net/http"
)

// ClientOption is a function type for configuring Client
type ClientOption func(*ClientConfig)

// WithPrivateKey sets the private key from App Store Connect (PEM format)
func WithPrivateKey(val []byte) ClientOption {
	return func(c *ClientConfig) {
		c.PrivateKey = val
	}
}

// WithKeyID sets the Key ID from App Store Connect
func WithKeyID(val string) ClientOption {
	return func(c *ClientConfig) {
		c.KeyID = val
	}
}

// WithIssuerID sets the Issuer ID from App Store Connect
func WithIssuerID(val string) ClientOption {
	return func(c *ClientConfig) {
		c.IssuerID = val
	}
}

// WithBundleID sets the Bundle ID of your app for Client
func WithBundleID(val string) ClientOption {
	return func(c *ClientConfig) {
		c.BundleID = val
	}
}

// WithEnvironment sets the environment (sandbox, production, etc.) for Client
func WithEnvironment(val Environment) ClientOption {
	return func(c *ClientConfig) {
		c.Environment = val
	}
}

// WithAppAppleID sets the App Apple ID (required for production environment) for Client
func WithAppAppleID(val int64) ClientOption {
	return func(c *ClientConfig) {
		c.AppAppleID = val
	}
}

// WithRootCertificates sets root certificates for JWS verification for Client
func WithRootCertificates(val [][]byte) ClientOption {
	return func(c *ClientConfig) {
		c.RootCertificates = val
	}
}

func WithEnableOnlineChecks() ClientOption {
	return func(config *ClientConfig) {
		config.EnableOnlineChecks = true
	}
}

// WithHTTPClient sets a custom HTTP client for API requests
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) {
		c.HTTPClient = client
	}
}

// WithEnableAutoDecode enables automatic decoding and verification of JWS in API responses
func WithEnableAutoDecode() ClientOption {
	return func(config *ClientConfig) {
		config.EnableAutoDecode = true
	}
}
