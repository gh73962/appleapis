package appstoreserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the App Store Server API
type APIError struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	HTTPStatus   int    `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("App Store Server API error: %d - %s (HTTP %d)", e.ErrorCode, e.ErrorMessage, e.HTTPStatus)
}

// NewAPIErrorFromResponse creates an APIError from an HTTP response
func NewAPIErrorFromResponse(resp *http.Response, body []byte) *APIError {
	var apiErr APIError
	apiErr.HTTPStatus = resp.StatusCode

	if len(body) > 0 {
		if err := json.Unmarshal(body, &apiErr); err != nil {
			apiErr.ErrorMessage = string(body)
		}
	} else {
		apiErr.ErrorMessage = resp.Status
	}

	return &apiErr
}

// VerificationError represents an error during verification
type VerificationError struct {
	Status VerificationStatus
	Err    error
}

// Error implements the error interface
func (e *VerificationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("verification failed with status %s: %v", e.Status.String(), e.Err)
	}
	return fmt.Sprintf("verification failed with status %s", e.Status.String())
}

// Unwrap returns the underlying error
func (e *VerificationError) Unwrap() error {
	return e.Err
}

// NewVerificationError creates a new verification error
func NewVerificationError(status VerificationStatus, err error) *VerificationError {
	return &VerificationError{
		Status: status,
		Err:    err,
	}
}
