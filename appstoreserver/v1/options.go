package appstoreserver

import "fmt"

// ClientConfig holds the configuration for creating a Client
type ClientConfig struct {
	PrivateKey         []byte
	KeyID              string
	IssuerID           string
	BundleID           string
	Environment        Environment
	AppAppleID         int64
	RootCertificates   [][]byte
	EnableOnlineChecks bool
	StrictChecks       bool
}

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

func WithStrictChecks() ClientOption {
	return func(config *ClientConfig) {
		config.StrictChecks = true
	}
}

// Validate validates the ClientConfig and returns an error if any required field is missing or invalid
func (c *ClientConfig) Validate() error {
	if len(c.PrivateKey) == 0 {
		return fmt.Errorf("private key is required")
	}
	if c.KeyID == "" {
		return fmt.Errorf("key ID is required")
	}
	if c.IssuerID == "" {
		return fmt.Errorf("issuer ID is required")
	}
	if c.BundleID == "" {
		return fmt.Errorf("bundle ID is required")
	}
	if !c.Environment.IsValid() {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}
	// Note: AppAppleID is optional for sandbox, required for production
	// This validation is handled in the New function context
	return nil
}
