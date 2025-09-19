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
	"strings"
	"testing"

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

// CreateTestClient creates a client configured for testing
func (h *TestHelper) CreateTestClient(t *testing.T) *Client {
	// Generate test private key in PEM format
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKGsJ1QwJ0WgKwZ5Mz3KZ5gEH8Z7c2+Y3+Ic7J8R8O9XoAoGCCqGSM49
AwEHoUQDQgAE4rWBxGmFbnPIPQI0zsBKzLxsj8pD2vqbr0yPISUx2WQyxmrNql9f
hK8YEEyYFV7++p5i4YUSR/o9uQIgCPIhrA==
-----END EC PRIVATE KEY-----`

	client, err := New(
		WithPrivateKey([]byte(privateKeyPEM)),
		WithKeyID("TESTKEY123"),
		WithIssuerID("TESTISSUER123"),
		WithBundleIDClient("com.example.test"),
		WithEnvironmentClient(EnvironmentLocalTesting),
		WithAppAppleIDClient(1234567890),
	)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	// Replace the base URL with our test server
	client.client.baseURL = h.server.URL

	return client
}

// ReadTestDataFile reads a test data file from the testdata directory
func ReadTestDataFile(relativePath string) ([]byte, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Navigate up to find the testdata directory
	for {
		testDataPath := filepath.Join(wd, "testdata", relativePath)
		if _, err := os.Stat(testDataPath); err == nil {
			return os.ReadFile(testDataPath)
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return nil, fmt.Errorf("test data file not found: %s", relativePath)
}

// ReadTestDataFileString reads a test data file and returns it as a string
func ReadTestDataFileString(relativePath string) (string, error) {
	data, err := ReadTestDataFile(relativePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CreateSignedDataFromJSON creates a JWT signed data from a JSON test file
func CreateSignedDataFromJSON(t *testing.T, filePath string) string {
	// Read the JSON data
	data, err := ReadTestDataFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read test data file %s: %v", filePath, err)
	}

	// Parse JSON to ensure it's valid
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("Failed to parse JSON from %s: %v", filePath, err)
	}

	// Generate a test private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(payload))

	// Add x5c header for certificate chain (using dummy cert)
	token.Header["x5c"] = []string{"dummycert"}

	// Sign the token
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign JWT: %v", err)
	}

	return signedToken
}

// AssertNoError is a test helper that fails the test if err is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// AssertError is a test helper that fails the test if err is nil
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

// AssertEqual compares two values and fails if they're not equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertStringEqual compares two strings and fails if they're not equal
func AssertStringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("Expected %q, got %q", expected, actual)
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Fatalf("Expected %q to contain %q", str, substr)
	}
}

// AssertNotNil checks if a value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Fatal("Expected non-nil value")
	}
}

// AssertNil checks if a value is nil
func AssertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		t.Fatalf("Expected nil, got %v", value)
	}
}

// AssertTrue checks if a value is true
func AssertTrue(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Fatal("Expected true")
	}
}

// AssertFalse checks if a value is false
func AssertFalse(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Fatal("Expected false")
	}
}
