package appstoreserver

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

// TestHelper provides utilities for testing the App Store Server API client
type TestHelper struct {
	privateKey *ecdsa.PrivateKey
	server     *httptest.Server
	responses  map[string]TestResponse
}

// TestResponse represents a mock HTTP response
type TestResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// NewTestHelper creates a new test helper
func NewTestHelper() *TestHelper {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	helper := &TestHelper{
		privateKey: privateKey,
		responses:  make(map[string]TestResponse),
	}

	// Create mock HTTP server
	helper.server = httptest.NewServer(http.HandlerFunc(helper.handleRequest))

	return helper
}

// Close cleans up the test helper
func (h *TestHelper) Close() {
	if h.server != nil {
		h.server.Close()
	}
}

// SetResponse sets a mock response for a specific request
func (h *TestHelper) SetResponse(method, path string, response TestResponse) {
	key := fmt.Sprintf("%s %s", method, path)
	h.responses[key] = response
}

// SetResponseFromFile sets a mock response from a test data file
func (h *TestHelper) SetResponseFromFile(method, path, filePath string) error {
	body, err := ReadTestDataFile(filePath)
	if err != nil {
		return err
	}

	h.SetResponse(method, path, TestResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	})

	return nil
}

// handleRequest handles incoming HTTP requests for the mock server
func (h *TestHelper) handleRequest(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

	response, exists := h.responses[key]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errorCode": 4040010, "errorMessage": "Not Found"}`))
		return
	}

	// Set headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}

// CreateSignedDataFromJSON creates a JWT signed data from a JSON test file
func CreateSignedDataFromJSON(filePath string) (string, error) {
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
