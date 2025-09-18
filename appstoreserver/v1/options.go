package appstoreserver

import "fmt"

// SignedDataVerifierConfig holds the configuration for creating a SignedDataVerifier
type SignedDataVerifierConfig struct {
	rootCertificates   [][]byte
	environment        Environment
	bundleID           string
	appAppleID         int64
	enableOnlineChecks bool
}

// SignedDataVerifierOption is a function type for configuring SignedDataVerifier
type SignedDataVerifierOption func(*SignedDataVerifierConfig)

// ClientConfig holds the configuration for creating a Client
type ClientConfig struct {
	PrivateKey       []byte
	KeyID            string
	IssuerID         string
	BundleID         string
	Environment      Environment
	AppAppleID       *int64
	RootCertificates [][]byte
}

// ClientOption is a function type for configuring Client
type ClientOption func(*ClientConfig)

// WithPrivateKey sets the private key from App Store Connect (PEM format)
func WithPrivateKey(privateKey []byte) ClientOption {
	return func(c *ClientConfig) {
		c.PrivateKey = privateKey
	}
}

// WithKeyID sets the Key ID from App Store Connect
func WithKeyID(keyID string) ClientOption {
	return func(c *ClientConfig) {
		c.KeyID = keyID
	}
}

// WithIssuerID sets the Issuer ID from App Store Connect
func WithIssuerID(issuerID string) ClientOption {
	return func(c *ClientConfig) {
		c.IssuerID = issuerID
	}
}

// WithBundleIDClient sets the Bundle ID of your app for Client
func WithBundleIDClient(bundleID string) ClientOption {
	return func(c *ClientConfig) {
		c.BundleID = bundleID
	}
}

// WithEnvironmentClient sets the environment (sandbox, production, etc.) for Client
func WithEnvironmentClient(environment Environment) ClientOption {
	return func(c *ClientConfig) {
		c.Environment = environment
	}
}

// WithAppAppleIDClient sets the App Apple ID (required for production environment) for Client
func WithAppAppleIDClient(appAppleID int64) ClientOption {
	return func(c *ClientConfig) {
		c.AppAppleID = &appAppleID
	}
}

// WithRootCertificatesClient sets root certificates for JWS verification for Client
func WithRootCertificatesClient(rootCertificates [][]byte) ClientOption {
	return func(c *ClientConfig) {
		c.RootCertificates = rootCertificates
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

// WithRootCertificates sets the root certificates for certificate chain verification
func WithRootCertificates(rootCertificates [][]byte) SignedDataVerifierOption {
	return func(config *SignedDataVerifierConfig) {
		config.rootCertificates = rootCertificates
	}
}

// WithEnvironment sets the App Store environment (Sandbox or Production)
func WithEnvironment(environment Environment) SignedDataVerifierOption {
	return func(config *SignedDataVerifierConfig) {
		config.environment = environment
	}
}

// WithBundleID sets the bundle ID to verify against
func WithBundleID(bundleID string) SignedDataVerifierOption {
	return func(config *SignedDataVerifierConfig) {
		config.bundleID = bundleID
	}
}

// WithAppAppleID sets the App Apple ID (required for Production environment)
func WithAppAppleID(appAppleID int64) SignedDataVerifierOption {
	return func(config *SignedDataVerifierConfig) {
		config.appAppleID = appAppleID
	}
}
