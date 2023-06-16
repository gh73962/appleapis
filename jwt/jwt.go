package jwt

import (
	"crypto/ecdsa"
	"os"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func NewToken(keyID string, claims *Claims, pk *ecdsa.PrivateKey) (*jwtv5.Token, string, error) {
	t := jwtv5.Token{
		Method: jwtv5.SigningMethodES256,
		Header: NewJWTHeader(keyID),
		Claims: claims,
	}

	bearer, err := t.SignedString(pk)
	if err != nil {
		return nil, "", err
	}

	return &t, bearer, nil
}

// GetPrivateKeyFromFile load apple xxxx.p8 certificate
func GetPrivateKeyFromFile(filename string) (*ecdsa.PrivateKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return jwtv5.ParseECPrivateKeyFromPEM(bytes)
}
