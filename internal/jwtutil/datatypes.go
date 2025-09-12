package jwtutil

import (
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
