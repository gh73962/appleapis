package appstoreserver

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gh73962/appleapis/appstoreservernotifications/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maypok86/otter/v2"
	"golang.org/x/crypto/ocsp"
)

const (
	// MaximumCacheSize limits the number of cached certificates
	MaximumCacheSize = 32
	// CacheTimeLimit defines how long certificates are cached (15 minutes)
	CacheTimeLimit = 15 * time.Minute
)

// chainVerifier handles certificate chain verification
type chainVerifier struct {
	rootCertificates   [][]byte
	cache              *otter.Cache[string, *ecdsa.PublicKey]
	enableStrictChecks bool
}

// generateCacheKey creates a unique cache key for a set of certificates
func generateCacheKey(certificates []string) string {
	if len(certificates) == 0 {
		return ""
	}
	// Use all certificates to create a unique key like Python version
	// which uses tuple(certificates) as key
	key := ""
	for _, cert := range certificates {
		key += cert + "|"
	}
	return key
}

// newChainVerifier creates a new chain verifier
func newChainVerifier(rootCertificates [][]byte) *chainVerifier {
	options := &otter.Options[string, *ecdsa.PublicKey]{
		MaximumSize:      MaximumCacheSize,
		ExpiryCalculator: otter.ExpiryWriting[string, *ecdsa.PublicKey](CacheTimeLimit),
	}

	cache, err := otter.New(options)
	if err != nil {
		panic(fmt.Sprintf("failed to create certificate cache: %v", err))
	}

	return &chainVerifier{
		rootCertificates: rootCertificates,
		cache:            cache,
	}
}

// verifyChain verifies the certificate chain and returns the public key for signature verification
func (c *chainVerifier) verifyChain(certificates []string, enableOnlineChecks bool, effectiveDate int64) (*ecdsa.PublicKey, error) {
	if enableOnlineChecks && len(certificates) > 0 {
		cacheKey := generateCacheKey(certificates)
		if cachedKey, found := c.cache.GetIfPresent(cacheKey); found {
			return cachedKey, nil
		}
	}

	publicKey, err := c.verifyChainWithoutCaching(certificates, enableOnlineChecks, effectiveDate)
	if err != nil {
		return nil, err
	}

	if enableOnlineChecks && len(certificates) > 0 {
		cacheKey := generateCacheKey(certificates)
		c.cache.Set(cacheKey, publicKey)
	}

	return publicKey, nil
}

// verifyChainWithoutCaching performs the actual certificate chain verification without caching
func (c *chainVerifier) verifyChainWithoutCaching(certificates []string, enableOnlineChecks bool, effectiveDate int64) (*ecdsa.PublicKey, error) {
	if len(c.rootCertificates) == 0 {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("no root certificates provided"))
	}

	if len(certificates) != 3 {
		return nil, NewVerificationError(VerificationStatusInvalidChainLength, fmt.Errorf("expected 3 certificates in chain, got %d", len(certificates)))
	}

	leafCertBytes, err := base64.StdEncoding.DecodeString(certificates[0])
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to decode leaf certificate: %w", err))
	}

	intermediateCertBytes, err := base64.StdEncoding.DecodeString(certificates[1])
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to decode intermediate certificate: %w", err))
	}

	rootCertBytes, err := base64.StdEncoding.DecodeString(certificates[2])
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to decode root certificate: %w", err))
	}

	leafCert, err := x509.ParseCertificate(leafCertBytes)
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to parse leaf certificate: %w", err))
	}

	intermediateCert, err := x509.ParseCertificate(intermediateCertBytes)
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to parse intermediate certificate: %w", err))
	}

	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, fmt.Errorf("failed to parse root certificate: %w", err))
	}

	if c.enableStrictChecks {
		var trusted bool
		for _, trustedRoot := range c.rootCertificates {
			trustedRootCert, err := x509.ParseCertificate(trustedRoot)
			if err != nil {
				continue
			}

			if rootCert.Equal(trustedRootCert) {
				trusted = true
				break
			}
		}

		if !trusted {
			return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("root certificate is not trusted"))
		}
	}

	roots := x509.NewCertPool()
	roots.AddCert(rootCert)

	intermediates := x509.NewCertPool()
	intermediates.AddCert(intermediateCert)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		CurrentTime:   time.Unix(effectiveDate, 0),
	}

	_, err = leafCert.Verify(opts)
	if err != nil {
		return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("certificate chain verification failed: %w", err))
	}

	if err := c.checkAppleOIDs(leafCert, intermediateCert); err != nil {
		return nil, err
	}

	if enableOnlineChecks {
		if err := c.checkOCSPStatus(leafCert, intermediateCert, rootCert); err != nil {
			return nil, err
		}
		if err := c.checkOCSPStatus(intermediateCert, rootCert, rootCert); err != nil {
			return nil, err
		}
	}

	publicKey, ok := leafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("leaf certificate does not contain ECDSA public key"))
	}

	return publicKey, nil
}

// checkAppleOIDs verifies that the certificates contain the required Apple OIDs
func (c *chainVerifier) checkAppleOIDs(leafCert, intermediateCert *x509.Certificate) error {
	leafOID := "1.2.840.113635.100.6.11.1"
	if !c.hasOID(leafCert, leafOID) {
		return NewVerificationError(VerificationStatusFailure, fmt.Errorf("leaf certificate missing required OID: %s", leafOID))
	}

	intermediateOID := "1.2.840.113635.100.6.2.1"
	if !c.hasOID(intermediateCert, intermediateOID) {
		return NewVerificationError(VerificationStatusFailure, fmt.Errorf("intermediate certificate missing required OID: %s", intermediateOID))
	}

	return nil
}

// hasOID checks if a certificate contains a specific OID extension
func (c *chainVerifier) hasOID(cert *x509.Certificate, oidStr string) bool {
	for _, ext := range cert.Extensions {
		if ext.Id.String() == oidStr {
			return true
		}
	}
	return false
}

// checkOCSPStatus performs OCSP (Online Certificate Status Protocol) checking
func (c *chainVerifier) checkOCSPStatus(cert, issuer, root *x509.Certificate) error {
	ocspRequest, err := ocsp.CreateRequest(cert, issuer, nil)
	if err != nil {
		return NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to create OCSP request: %w", err))
	}

	ocspServerURLs := c.getOCSPServerURLs(cert)
	if len(ocspServerURLs) == 0 {
		return NewVerificationError(VerificationStatusFailure, errors.New("no OCSP server URLs found in certificate"))
	}

	for _, serverURL := range ocspServerURLs {
		response, err := c.queryOCSPServer(serverURL, ocspRequest)
		if err != nil {
			continue
		}

		ocspResponse, err := ocsp.ParseResponse(response, nil)
		if err != nil {
			continue
		}

		if ocspResponse.Status == ocsp.Good {
			// Additional verification like Python version
			// Re-parse with issuer for proper signature verification
			ocspResponseVerified, err := ocsp.ParseResponse(response, issuer)
			if err != nil {
				continue
			}
			if ocspResponseVerified.Status == ocsp.Good {
				return nil
			}
		}

		return NewVerificationError(VerificationStatusFailure, fmt.Errorf("certificate status is not good: %d", ocspResponse.Status))
	}

	return NewVerificationError(VerificationStatusFailure, errors.New("failed to verify certificate status via OCSP"))
}

// getOCSPServerURLs extracts OCSP server URLs from certificate's AIA extension
func (c *chainVerifier) getOCSPServerURLs(cert *x509.Certificate) []string {
	var ocspURLs []string

	for _, ext := range cert.Extensions {
		if ext.Id.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 1}) {
			var aia []struct {
				Method   asn1.ObjectIdentifier
				Location asn1.RawValue
			}

			if _, err := asn1.Unmarshal(ext.Value, &aia); err == nil {
				for _, access := range aia {
					if access.Method.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1}) {
						if access.Location.Tag == 6 {
							ocspURLs = append(ocspURLs, string(access.Location.Bytes))
						}
					}
				}
			}
		}
	}

	return ocspURLs
}

// queryOCSPServer sends OCSP request to server and returns response
func (c *chainVerifier) queryOCSPServer(serverURL string, request []byte) ([]byte, error) {
	_, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid OCSP server URL: %w", err)
	}

	resp, err := http.Post(serverURL, "application/ocsp-request", bytes.NewReader(request))
	if err != nil {
		return nil, fmt.Errorf("failed to send OCSP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCSP server returned status: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OCSP response: %w", err)
	}

	return responseBody, nil
}

// SignedDataVerifier provides utility methods for verifying and decoding App Store signed data
type SignedDataVerifier struct {
	rootCertificates   [][]byte
	environment        Environment
	bundleID           string
	appAppleID         int64
	chainVerifier      *chainVerifier
	enableOnlineChecks bool
}

// NewSignedDataVerifier creates a new SignedDataVerifier instance
func NewSignedDataVerifier(rootCertificates [][]byte, enableOnlineChecks bool, environment Environment, bundleID string, appAppleID int64) (*SignedDataVerifier, error) {
	if environment == EnvironmentProduction && appAppleID == 0 {
		return nil, errors.New("appAppleID is required when the environment is Production")
	}

	return &SignedDataVerifier{
		rootCertificates:   rootCertificates,
		environment:        environment,
		bundleID:           bundleID,
		appAppleID:         appAppleID,
		chainVerifier:      newChainVerifier(rootCertificates),
		enableOnlineChecks: enableOnlineChecks,
	}, nil
}

// VerifyAndDecodeRenewalInfo verifies and decodes a signedRenewalInfo obtained from the App Store Server API
func (v *SignedDataVerifier) VerifyAndDecodeRenewalInfo(signedRenewalInfo string) (*JWSRenewalInfoDecodedPayload, error) {
	decodedPayload, err := v.decodeSignedObject(signedRenewalInfo)
	if err != nil {
		return nil, err
	}

	var renewalInfo JWSRenewalInfoDecodedPayload
	if err := json.Unmarshal(decodedPayload, &renewalInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal renewal info: %w", err)
	}

	if renewalInfo.Environment != v.environment {
		return nil, NewVerificationError(VerificationStatusInvalidEnvironment, nil)
	}

	return &renewalInfo, nil
}

// VerifyAndDecodeSignedTransaction verifies and decodes a signedTransaction obtained from the App Store Server API
func (v *SignedDataVerifier) VerifyAndDecodeSignedTransaction(signedTransaction string) (*JWSTransactionDecodedPayload, error) {
	decodedPayload, err := v.decodeSignedObject(signedTransaction)
	if err != nil {
		return nil, err
	}

	var transactionInfo JWSTransactionDecodedPayload
	if err := json.Unmarshal(decodedPayload, &transactionInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction info: %w", err)
	}

	if transactionInfo.BundleID != v.bundleID {
		return nil, NewVerificationError(VerificationStatusInvalidAppIdentifier, nil)
	}

	if transactionInfo.Environment != v.environment {
		return nil, NewVerificationError(VerificationStatusInvalidEnvironment, nil)
	}

	return &transactionInfo, nil
}

// VerifyAndDecodeNotification verifies and decodes an App Store Server Notification signedPayload
func (v *SignedDataVerifier) VerifyAndDecodeNotification(signedPayload string) (*appstoreservernotifications.ResponseBodyV2DecodedPayload, error) {
	decodedPayload, err := v.decodeSignedObject(signedPayload)
	if err != nil {
		return nil, err
	}

	var notification appstoreservernotifications.ResponseBodyV2DecodedPayload
	if err := json.Unmarshal(decodedPayload, &notification); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	var (
		bundleID    string
		appAppleID  int64
		environment string
	)

	switch {
	case notification.Data != nil:
		bundleID = notification.Data.BundleID
		appAppleID = notification.Data.AppAppleID
		environment = notification.Data.Environment

	case notification.Summary != nil:
		bundleID = notification.Summary.BundleID
		appAppleID = notification.Summary.AppAppleID

	case notification.ExternalPurchaseToken != nil:
		bundleID = notification.ExternalPurchaseToken.BundleID
		appAppleID = notification.ExternalPurchaseToken.AppAppleID
		if strings.HasPrefix(notification.ExternalPurchaseToken.ExternalPurchaseID, "SANDBOX") {
			environment = EnvironmentSandbox.String()
		} else {
			environment = EnvironmentProduction.String()
		}
	}

	if bundleID != v.bundleID || (v.environment == EnvironmentProduction && appAppleID != v.appAppleID) {
		return nil, NewVerificationError(VerificationStatusInvalidAppIdentifier, nil)
	}

	if Environment(environment) != v.environment {
		return nil, NewVerificationError(VerificationStatusInvalidEnvironment, nil)
	}

	return &notification, nil
}

// decodeSignedObject decodes and verifies a signed JWT object
func (v *SignedDataVerifier) decodeSignedObject(signedObj string) ([]byte, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(signedObj, jwt.MapClaims{})
	if err != nil {
		return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to parse JWT: %w", err))
	}

	if v.environment == EnvironmentLocalTesting {
		claimsBytes, err := json.Marshal(token.Claims)
		if err != nil {
			return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to marshal claims: %w", err))
		}
		return claimsBytes, nil
	}

	x5cHeader, ok := token.Header["x5c"].([]any)
	if !ok || len(x5cHeader) == 0 {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("x5c claim was empty"))
	}

	alg, ok := token.Header["alg"].(string)
	if !ok || alg != "ES256" {
		return nil, NewVerificationError(VerificationStatusFailure, errors.New("algorithm was not ES256"))
	}

	certificates := make([]string, len(x5cHeader))
	for i, cert := range x5cHeader {
		certStr, ok := cert.(string)
		if !ok {
			return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("invalid certificate in x5c header"))
		}
		certificates[i] = certStr
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewVerificationError(VerificationStatusFailure, errors.New("failed to get claims"))
	}

	var effectiveDate int64
	if v.enableOnlineChecks {
		effectiveDate = time.Now().Unix()
	} else {
		if signedDate, exists := claims["signedDate"]; exists {
			if signedDateFloat, ok := signedDate.(float64); ok {
				effectiveDate = int64(signedDateFloat) / 1000
			} else {
				effectiveDate = time.Now().Unix()
			}
		} else if receiptCreationDate, exists := claims["receiptCreationDate"]; exists {
			if receiptCreationDateFloat, ok := receiptCreationDate.(float64); ok {
				effectiveDate = int64(receiptCreationDateFloat) / 1000
			} else {
				effectiveDate = time.Now().Unix()
			}
		} else {
			effectiveDate = time.Now().Unix()
		}
	}

	signingKey, err := v.chainVerifier.verifyChain(certificates, v.enableOnlineChecks, effectiveDate)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(signedObj, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to verify JWT signature: %w", err))
	}

	if !parsedToken.Valid {
		return nil, NewVerificationError(VerificationStatusFailure, errors.New("JWT token is invalid"))
	}

	claimsBytes, err := json.Marshal(parsedToken.Claims)
	if err != nil {
		return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to marshal verified claims: %w", err))
	}

	return claimsBytes, nil
}
