package appstoreserver

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
