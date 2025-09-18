package appstoreserver

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims see https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
type Claims struct {
	Issuer         string `json:"iss"`
	IssuedAt       int64  `json:"iat"`
	ExpirationTime int64  `json:"exp"`
	Audience       string `json:"aud"` // Audience must appstoreconnect-v1
	BundleID       string `json:"bid"`
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.ExpirationTime, 0)}, nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.IssuedAt, 0)}, nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return []string{c.Audience}, nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c *Claims) GetSubject() (string, error) {
	return "", nil
}

func NewClaims(iss, bid string) *Claims {
	t := time.Now()
	return &Claims{
		Issuer:         iss,
		IssuedAt:       t.Unix(),
		ExpirationTime: t.Add(30 * time.Minute).Unix(),
		Audience:       "appstoreconnect-v1",
		BundleID:       bid,
	}
}

func NewJWTHeader(keyID string) map[string]any {
	return map[string]any{
		"alg": "ES256",
		"kid": keyID,
		"typ": "JWT",
	}
}

const (
	// JWTAudience is the required audience for App Store Server API JWT tokens
	JWTAudience = "appstoreconnect-v1"
)

// TokenGenerator generates JWT tokens for App Store Server API authentication
type TokenGenerator struct {
	signingKey *ecdsa.PrivateKey
	keyID      string
	issuerID   string
	bundleID   string
}

// NewTokenGenerator creates a new JWT token generator
func NewTokenGenerator(signingKey *ecdsa.PrivateKey, keyID, issuerID, bundleID string) *TokenGenerator {
	return &TokenGenerator{
		signingKey: signingKey,
		keyID:      keyID,
		issuerID:   issuerID,
		bundleID:   bundleID,
	}
}

// GenerateToken creates a new JWT token for API authentication
func (t *TokenGenerator) GenerateToken() (string, error) {
	now := time.Now()
	claims := &Claims{
		Issuer:         t.issuerID,
		IssuedAt:       now.Unix(),
		ExpirationTime: now.Add(5 * time.Minute).Unix(), // Token expires in 5 minutes
		Audience:       JWTAudience,
		BundleID:       t.bundleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header = NewJWTHeader(t.keyID)

	tokenString, err := token.SignedString(t.signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ParsePrivateKeyFromPEM parses an ECDSA private key from PEM format
func ParsePrivateKeyFromPEM(pemData []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "PRIVATE KEY":
		// PKCS#8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
		}

		ecdsaKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not an ECDSA key")
		}

		return ecdsaKey, nil

	case "EC PRIVATE KEY":
		// SEC 1 format
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse EC private key: %w", err)
		}

		return key, nil

	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}
