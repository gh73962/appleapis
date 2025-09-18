package appstoreserver

import (
	"fmt"
)

// Config contains configuration for the App Store Server SDK
type Config struct {
	// Private key from App Store Connect (PEM format)
	PrivateKey []byte
	// Key ID from App Store Connect
	KeyID string
	// Issuer ID from App Store Connect
	IssuerID string
	// Bundle ID of your app
	BundleID string
	// Environment (sandbox, production, etc.)
	Environment Environment
	// App Apple ID (required for production environment)
	AppAppleID *int64
	// Root certificates for JWS verification (optional, uses Apple's by default)
	RootCertificates [][]byte
	// Enable online checks for certificate validation
	EnableOnlineChecks bool
}

// SDK provides a high-level interface to the App Store Server API and JWS verification
type SDK struct {
	client   *Client
	verifier *SignedDataVerifier
}

// NewSDK creates a new App Store Server SDK instance
func NewSDK(config Config) (*SDK, error) {
	// Validate required fields
	if len(config.PrivateKey) == 0 {
		return nil, fmt.Errorf("private key is required")
	}
	if config.KeyID == "" {
		return nil, fmt.Errorf("key ID is required")
	}
	if config.IssuerID == "" {
		return nil, fmt.Errorf("issuer ID is required")
	}
	if config.BundleID == "" {
		return nil, fmt.Errorf("bundle ID is required")
	}
	if !config.Environment.IsValid() {
		return nil, fmt.Errorf("invalid environment: %s", config.Environment)
	}

	// Parse private key
	privateKey, err := ParsePrivateKeyFromPEM(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create token generator
	tokenGenerator := NewTokenGenerator(privateKey, config.KeyID, config.IssuerID, config.BundleID)

	// Create API client
	client, err := NewClient(config.Environment, tokenGenerator)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Create JWS verifier
	var rootCerts [][]byte = config.RootCertificates
	// TODO: If no root certificates provided, use Apple's default root certificates

	options := []SignedDataVerifierOption{
		WithRootCertificates(rootCerts),
		WithEnvironment(config.Environment),
		WithBundleID(config.BundleID),
	}

	// Add AppAppleID only if it's provided (not nil)
	if config.AppAppleID != nil {
		options = append(options, WithAppAppleID(*config.AppAppleID))
	}

	verifier, err := NewSignedDataVerifier(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWS verifier: %w", err)
	}

	return &SDK{
		client:   client,
		verifier: verifier,
	}, nil
}

// GetClient returns the underlying API client
func (s *SDK) GetClient() *Client {
	return s.client
}

// GetVerifier returns the underlying JWS verifier
func (s *SDK) GetVerifier() *SignedDataVerifier {
	return s.verifier
}

// VerifyAndDecodeTransaction verifies and decodes a signed transaction
func (s *SDK) VerifyAndDecodeTransaction(signedTransaction string) (*JWSTransactionDecodedPayload, error) {
	return s.verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
}

// VerifyAndDecodeRenewalInfo verifies and decodes signed renewal info
func (s *SDK) VerifyAndDecodeRenewalInfo(signedRenewalInfo string) (*JWSRenewalInfoDecodedPayload, error) {
	return s.verifier.VerifyAndDecodeRenewalInfo(signedRenewalInfo)
}
