package appstoreserver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientConfig holds the configuration for creating a Client
// ClientConfig contains the configuration parameters required to initialize
// an App Store Server API client for communicating with Apple's services.
type ClientConfig struct {
	// PrivateKey is the private key content in PEM format used for JWT signing.
	// This key should correspond to the key downloaded from App Store Connect.
	PrivateKey []byte

	// KeyID is the key identifier from App Store Connect, used to identify
	// which key was used to sign the JWT token.
	KeyID string

	// IssuerID is the issuer identifier from App Store Connect, representing
	// your developer account or organization.
	IssuerID string

	// BundleID is the app's bundle identifier (e.g., com.example.myapp)
	// used to identify the specific app for API requests.
	BundleID string

	// Environment specifies whether to use the sandbox or production
	// App Store Server API environment.
	Environment Environment

	// AppAppleID is the unique identifier assigned by Apple to your app,
	// found in App Store Connect under App Information.
	AppAppleID int64

	// RootCertificates contains the Apple Root CA certificates in DER format
	// used for verifying signed data from Apple's servers.
	// defaults to well-known Apple Root CAs if not provided.
	RootCertificates [][]byte

	// EnableOnlineChecks determines whether to perform online verification
	// of certificates and CRL (Certificate Revocation List) checking.
	EnableOnlineChecks bool

	// HTTPClient is the custom HTTP client to use for API requests.
	// If nil, a default HTTP client will be used.
	HTTPClient *http.Client

	// EnableAutoDecode controls whether responses should be automatically
	// decoded and verified. When true, JWS signatures are verified automatically.
	// recommended to be true for most use cases.
	// When false, raw responses are returned and must be verified manually.
	EnableAutoDecode bool
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
	if c.Environment == EnvironmentProduction && c.AppAppleID == 0 {
		return errors.New("appAppleID is required when the environment is Production")
	}
	return nil
}

func (c *ClientConfig) Init() error {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if len(c.RootCertificates) == 0 {
		for _, v := range []string{AppleRootCAURL, AppleRootCAG2URL, AppleRootCAG3URL} {
			err := func() error {
				resp, err := http.Get(v)
				if err != nil {
					return fmt.Errorf("failed to download %s: %w", v, err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("failed to download %s: HTTP %d", v, resp.StatusCode)
				}

				certData, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to read certificate data from %s: %w", v, err)
				}

				c.RootCertificates = append(c.RootCertificates, certData)
				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
	if c.Environment == "" {
		c.Environment = EnvironmentSandbox
	}
	return nil
}
