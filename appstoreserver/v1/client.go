package appstoreserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client represents the App Store Server API client
type Client struct {
	baseURL        string
	tokenGenerator *TokenGenerator
	httpClient     *http.Client
	userAgent      string
}

// NewClient creates a new App Store Server API client
func NewClient(environment Environment, tokenGenerator *TokenGenerator) (*Client, error) {
	var baseURL string
	switch environment {
	case EnvironmentProduction:
		baseURL = ProductionBaseURL
	case EnvironmentSandbox:
		baseURL = SandboxBaseURL
	default:
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}

	return &Client{
		baseURL:        baseURL,
		tokenGenerator: tokenGenerator,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "app-store-server-library/go/1.0.0",
	}, nil
}

// GetTransactionInfo gets information about a single transaction
// See https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (c *Client) GetTransactionInfo(ctx context.Context, transactionID string) (*TransactionInfoResponse, error) {
	path := fmt.Sprintf("/inApps/v1/transactions/%s", transactionID)

	var response TransactionInfoResponse
	if err := c.makeRequest(ctx, "GET", path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTransactionHistory gets a customer's transaction history for your app
// See https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (c *Client) GetTransactionHistory(ctx context.Context, transactionID string, request *TransactionHistoryRequest) (*HistoryResponse, error) {
	path := fmt.Sprintf("/inApps/v2/history/%s", transactionID)

	queryParams := make(url.Values)
	if request != nil {
		if request.Sort != "" {
			queryParams.Add("sort", request.Sort)
		}
		if len(request.ProductIDs) > 0 {
			for _, productID := range request.ProductIDs {
				queryParams.Add("productId", productID)
			}
		}
		if len(request.ProductTypes) > 0 {
			for _, productType := range request.ProductTypes {
				queryParams.Add("productType", productType)
			}
		}
		if request.StartDate != 0 {
			queryParams.Add("startDate", strconv.FormatInt(request.StartDate, 10))
		}
		if request.EndDate != 0 {
			queryParams.Add("endDate", strconv.FormatInt(request.EndDate, 10))
		}
		if len(request.SubscriptionGroupIdentifiers) > 0 {
			for _, id := range request.SubscriptionGroupIdentifiers {
				queryParams.Add("subscriptionGroupIdentifier", id)
			}
		}
		if request.InAppOwnershipType != "" {
			queryParams.Add("inAppOwnershipType", request.InAppOwnershipType)
		}
		if request.Revoked {
			queryParams.Add("revoked", strconv.FormatBool(request.Revoked))
		}
	}

	var response HistoryResponse
	if err := c.makeRequest(ctx, "GET", path, queryParams, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// LookUpOrderID gets a customer's in-app purchases from a receipt using the order ID
// See https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id
func (c *Client) LookUpOrderID(ctx context.Context, orderID string) (*OrderLookupResponse, error) {
	path := fmt.Sprintf("/inApps/v1/lookup/%s", orderID)

	var response OrderLookupResponse
	if err := c.makeRequest(ctx, "GET", path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetRefundHistory gets a paginated list of all of a customer's refunded in-app purchases
// See https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (c *Client) GetRefundHistory(ctx context.Context, transactionID string, revision *string) (*RefundHistoryResponse, error) {
	path := fmt.Sprintf("/inApps/v2/refund/lookup/%s", transactionID)

	queryParams := make(url.Values)
	if revision != nil {
		queryParams.Add("revision", *revision)
	}

	var response RefundHistoryResponse
	if err := c.makeRequest(ctx, "GET", path, queryParams, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ExtendRenewalDate extends the renewal date for a subscription
// See https://developer.apple.com/documentation/appstoreserverapi/extend_a_subscription_renewal_date
func (c *Client) ExtendRenewalDate(ctx context.Context, originalTransactionID string, request *ExtendRenewalDateRequest) (*ExtendRenewalDateResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/%s", originalTransactionID)

	var response ExtendRenewalDateResponse
	if err := c.makeRequest(ctx, "PUT", path, nil, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// MassExtendRenewalDate extends the renewal date for all active subscribers
// See https://developer.apple.com/documentation/appstoreserverapi/extend_subscription_renewal_dates_for_all_active_subscribers
func (c *Client) MassExtendRenewalDate(ctx context.Context, request *MassExtendRenewalDateRequest) (*MassExtendRenewalDateResponse, error) {
	path := "/inApps/v1/subscriptions/extend/mass"

	var response MassExtendRenewalDateResponse
	if err := c.makeRequest(ctx, "POST", path, nil, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMassExtendRenewalDateStatus checks the status of a mass renewal date extension request
// See https://developer.apple.com/documentation/appstoreserverapi/get_status_of_subscription_renewal_date_extensions
func (c *Client) GetMassExtendRenewalDateStatus(ctx context.Context, productID, requestIdentifier string) (*MassExtendRenewalDateStatusResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/mass/%s/%s", productID, requestIdentifier)

	var response MassExtendRenewalDateStatusResponse
	if err := c.makeRequest(ctx, "GET", path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetNotificationHistory gets a list of notifications that the App Store server attempted to send
// See https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
func (c *Client) GetNotificationHistory(ctx context.Context, paginationToken *string, request *NotificationHistoryRequest) (*NotificationHistoryResponse, error) {
	path := "/inApps/v1/notifications/history"

	queryParams := make(url.Values)
	if paginationToken != nil {
		queryParams.Add("paginationToken", *paginationToken)
	}

	var response NotificationHistoryResponse
	if err := c.makeRequest(ctx, "POST", path, queryParams, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// RequestTestNotification asks App Store Server Notifications to send a test notification
// See https://developer.apple.com/documentation/appstoreserverapi/request_a_test_notification
func (c *Client) RequestTestNotification(ctx context.Context) (*SendTestNotificationResponse, error) {
	path := "/inApps/v1/notifications/test"

	var response SendTestNotificationResponse
	if err := c.makeRequest(ctx, "POST", path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTestNotificationStatus checks the status of a test notification
// See https://developer.apple.com/documentation/appstoreserverapi/get_test_notification_status
func (c *Client) GetTestNotificationStatus(ctx context.Context, testNotificationToken string) (*CheckTestNotificationResponse, error) {
	path := fmt.Sprintf("/inApps/v1/notifications/test/%s", testNotificationToken)

	var response CheckTestNotificationResponse
	if err := c.makeRequest(ctx, "GET", path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendConsumptionData sends consumption information about a consumable in-app purchase
// See https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information
func (c *Client) SendConsumptionData(ctx context.Context, transactionID string, request *ConsumptionRequest) error {
	path := fmt.Sprintf("/inApps/v1/transactions/consumption/%s", transactionID)

	return c.makeRequest(ctx, "PUT", path, nil, request, nil)
}

// SetAppAccountToken sets the app account token value for a purchase
// See https://developer.apple.com/documentation/appstoreserverapi/set-app-account-token
func (c *Client) SetAppAccountToken(ctx context.Context, originalTransactionID string, request *UpdateAppAccountTokenRequest) error {
	path := fmt.Sprintf("/inApps/v1/transactions/%s/appAccountToken", originalTransactionID)

	return c.makeRequest(ctx, "PUT", path, nil, request, nil)
}

// makeRequest performs an HTTP request to the App Store Server API
func (c *Client) makeRequest(ctx context.Context, method, path string, queryParams url.Values, requestBody, responseBody any) error {
	// Generate JWT token
	token, err := c.tokenGenerator.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Build URL
	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	// Prepare request body
	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewAPIErrorFromResponse(resp, respBodyBytes)
	}

	// Parse response body if needed
	if responseBody != nil {
		if err := json.Unmarshal(respBodyBytes, responseBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
