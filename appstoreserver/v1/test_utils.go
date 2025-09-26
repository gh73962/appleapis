package appstoreserver

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

func mockSignedData(filePath string) (string, error) {
	data, err := os.ReadFile(filepath.Join("../../testdata/", filePath))
	if err != nil {
		return "", err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(payload))

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func mockTestClient(opts ...ClientOption) (*client, error) {
	pk, err := os.ReadFile("../../testdata/certs/testSigningKey.p8")
	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile("../../testdata/certs/testCA.der")
	if err != nil {
		return nil, err
	}

	testOps := []ClientOption{
		WithAppAppleID(1234),
		WithBundleID("com.example"),
		WithEnvironment(EnvironmentLocalTesting),
		WithKeyID("keyId"),
		WithIssuerID("issuerId"),
		WithPrivateKey(pk),
		WithRootCertificates([][]byte{cert}),
	}

	testOps = append(testOps, opts...)

	client, err := New(testOps...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func mockClientWithBody(filePath string, statusCode int, opts ...ClientOption) (*client, error) {
	var (
		responseBody []byte
		err          error
	)
	if filePath != "" {
		responseBody, err = os.ReadFile(filepath.Join("../../testdata/", filePath))
		if err != nil {
			return nil, err
		}
	}

	mockTransport := &mockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			}, nil
		},
	}

	httpClient := &http.Client{
		Transport: mockTransport,
	}

	opts = append(opts, WithHTTPClient(httpClient))
	return mockTestClient(opts...)
}
