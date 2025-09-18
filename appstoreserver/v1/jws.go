package appstoreserver

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gh73962/appleapis/appstoreservernotifications/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// MaximumCacheSize limits the number of cached certificates
	MaximumCacheSize = 32
	// CacheTimeLimit defines how long certificates are cached (15 minutes)
	CacheTimeLimit = 15 * 60 * time.Second
)

// chainVerifier handles certificate chain verification
type chainVerifier struct {
	rootCertificates [][]byte
	// For now, we'll implement a basic version without caching and OCSP checking
	// This can be extended later with more advanced features
}

// newChainVerifier creates a new chain verifier
func newChainVerifier(rootCertificates [][]byte) *chainVerifier {
	return &chainVerifier{
		rootCertificates: rootCertificates,
	}
}

// verifyChain verifies the certificate chain and returns the public key for signature verification
func (c *chainVerifier) verifyChain(certificates []string, effectiveDate int64) (*ecdsa.PublicKey, error) {
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

	trusted := false
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

	publicKey, ok := leafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, NewVerificationError(VerificationStatusInvalidCertificate, errors.New("leaf certificate does not contain ECDSA public key"))
	}

	// TODO: Implement OCSP checking if performOnlineChecks is true
	// This would require additional dependencies and complexity

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

// SignedDataVerifier provides utility methods for verifying and decoding App Store signed data
type SignedDataVerifier struct {
	rootCertificates   [][]byte
	environment        Environment
	bundleID           string
	appAppleID         int64
	chainVerifier      *chainVerifier
	enableOnlineChecks bool
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
		bundleID     string
		appAppleID   int64
		environment  string
		hasValidData bool
	)

	switch {
	case notification.Data != nil:
		bundleID = notification.Data.BundleID
		appAppleID = notification.Data.AppAppleID
		environment = notification.Data.Environment
		hasValidData = true
	case notification.Summary != nil:
		bundleID = notification.Summary.BundleID
		appAppleID = notification.Summary.AppAppleID
		environment = notification.Summary.Environment
		hasValidData = true
	case notification.ExternalPurchaseToken != nil:
		bundleID = notification.ExternalPurchaseToken.BundleID
		appAppleID = notification.ExternalPurchaseToken.AppAppleID
		// Determine environment based on externalPurchaseId
		if len(notification.ExternalPurchaseToken.ExternalPurchaseID) > 7 &&
			notification.ExternalPurchaseToken.ExternalPurchaseID[:7] == "SANDBOX" {
			environment = EnvironmentSandbox.String()
		} else {
			environment = EnvironmentProduction.String()
		}
		hasValidData = true
	}

	if !hasValidData {
		return nil, NewVerificationError(VerificationStatusFailure, errors.New("notification contains no valid data"))
	}

	if err := v.verifyNotification(bundleID, appAppleID, Environment(environment)); err != nil {
		return nil, err
	}

	return &notification, nil
}

// verifyNotification validates the notification's bundleID, appAppleID, and environment
func (v *SignedDataVerifier) verifyNotification(bundleID string, appAppleID int64, environment Environment) error {
	if bundleID != v.bundleID {
		return NewVerificationError(VerificationStatusInvalidAppIdentifier, nil)
	}

	if v.environment == EnvironmentProduction && appAppleID != v.appAppleID {
		return NewVerificationError(VerificationStatusInvalidAppIdentifier, nil)

	}

	if environment != v.environment {
		return NewVerificationError(VerificationStatusInvalidEnvironment, nil)
	}

	return nil
}

// decodeSignedObject decodes and verifies a signed JWT object
func (v *SignedDataVerifier) decodeSignedObject(signedObj string) ([]byte, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(signedObj, jwt.MapClaims{})
	if err != nil {
		return nil, NewVerificationError(VerificationStatusFailure, fmt.Errorf("failed to parse JWT: %w", err))
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

	signingKey, err := v.chainVerifier.verifyChain(certificates, effectiveDate)
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
