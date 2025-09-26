package appstoreserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetTransactionHistory gets a customer's transaction history for your app
// See https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (c *Client) GetTransactionHistory(ctx context.Context, req *TransactionHistoryRequest) (*HistoryResponse, error) {
	if req.TransactionID == "" {
		return nil, fmt.Errorf("transactionID cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v2/history/%s", req.TransactionID)

	queryParams := req.makeQuery()

	var response HistoryResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, queryParams, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}
	for _, v := range response.SignedTransactions {
		payload, err := c.Verifier.VerifyAndDecodeSignedTransaction(v)
		if err != nil {
			return nil, fmt.Errorf("SignedTransactions %s\nfailed to verify and decode: %w", v, err)
		}
		response.Payloads = append(response.Payloads, payload)
	}

	return &response, nil
}

// GetTransactionInfo gets information about a single transaction
// See https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (c *Client) GetTransactionInfo(ctx context.Context, transactionID string) (*TransactionInfoResponse, error) {
	if transactionID == "" {
		return nil, fmt.Errorf("transactionID cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v1/transactions/%s", transactionID)

	var response TransactionInfoResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	payload, err := c.Verifier.VerifyAndDecodeSignedTransaction(response.SignedTransactionInfo)
	if err != nil {
		return nil, fmt.Errorf("SignedTransactions %s\nfailed to verify and decode: %w", response.SignedTransactionInfo, err)
	}
	response.Payload = payload

	return &response, nil
}

// GetAllSubscriptionStatuses gets the statuses for all of a customer's auto-renewable subscriptions in your app
// See https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
func (c *Client) GetAllSubscriptionStatuses(ctx context.Context, transactionID string, status ...SubscriptionStatus) (*StatusResponse, error) {
	if transactionID == "" {
		return nil, fmt.Errorf("transactionID cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v1/subscriptions/%s", transactionID)

	queryParams := make(url.Values)
	if len(status) > 0 {
		for _, s := range status {
			queryParams.Add("status", s.String())
		}
	}

	var response StatusResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, queryParams, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	for _, data := range response.Data {
		for i, v := range data.LastTransactions {
			renewalPayload, err := c.Verifier.VerifyAndDecodeRenewalInfo(v.SignedRenewalInfo)
			if err != nil {
				return nil, fmt.Errorf("SignedRenewalInfo %s\nfailed to verify and decode: %w", v.SignedRenewalInfo, err)
			}
			data.LastTransactions[i].RenewalPayload = renewalPayload
			transactionPayload, err := c.Verifier.VerifyAndDecodeSignedTransaction(v.SignedTransactionInfo)
			if err != nil {
				return nil, fmt.Errorf("SignedTransactionInfo %s\nfailed to verify and decode: %w", v.SignedTransactionInfo, err)
			}
			data.LastTransactions[i].TransactionPayload = transactionPayload
		}
	}

	return &response, nil
}

// SetAppAccountToken sets the app account token value for a purchase
// See https://developer.apple.com/documentation/appstoreserverapi/set-app-account-token
func (c *Client) SetAppAccountToken(ctx context.Context, req *UpdateAppAccountTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/inApps/v1/transactions/%s/appAccountToken", req.OriginalTransactionID)
	return c.makeRequest(ctx, http.MethodPut, path, nil, req, nil)
}

// SendConsumptionData sends consumption information about a consumable in-app purchase
// See https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information
func (c *Client) SendConsumptionInfo(ctx context.Context, req *ConsumptionRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/inApps/v1/transactions/consumption/%s", req.TransactionID)
	return c.makeRequest(ctx, http.MethodPut, path, nil, req, nil)
}

// LookUpOrderID gets a customer's in-app purchases from a receipt using the order ID
// See https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id
func (c *Client) LookUpOrderID(ctx context.Context, orderID string) (*OrderLookupResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("orderID cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v1/lookup/%s", orderID)

	var response OrderLookupResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	for _, v := range response.SignedTransactions {
		payload, err := c.Verifier.VerifyAndDecodeSignedTransaction(v)
		if err != nil {
			return nil, fmt.Errorf("SignedTransactions %s\nfailed to verify and decode: %w", v, err)
		}
		response.Payloads = append(response.Payloads, payload)
	}

	return &response, nil
}

// GetRefundHistory gets a paginated list of all of a customer's refunded in-app purchases
// See https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (c *Client) GetRefundHistory(ctx context.Context, transactionID, revision string) (*RefundHistoryResponse, error) {
	if transactionID == "" {
		return nil, fmt.Errorf("transactionID cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v2/refund/lookup/%s", transactionID)

	queryParams := make(url.Values)
	if revision != "" {
		queryParams.Add("revision", revision)
	}

	var response RefundHistoryResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, queryParams, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	for _, v := range response.SignedTransactions {
		payload, err := c.Verifier.VerifyAndDecodeSignedTransaction(v)
		if err != nil {
			return nil, fmt.Errorf("SignedTransactions %s\nfailed to verify and decode: %w", v, err)
		}
		response.Payloads = append(response.Payloads, payload)
	}

	return &response, nil
}

// ExtendSubscriptionRenewalDate extends the renewal date for a subscription
// See https://developer.apple.com/documentation/appstoreserverapi/extend_a_subscription_renewal_date
func (c *Client) ExtendSubscriptionRenewalDate(ctx context.Context, req *ExtendRenewalDateRequest) (*ExtendRenewalDateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/%s", req.OriginalTransactionID)

	var response ExtendRenewalDateResponse
	if err := c.makeRequest(ctx, http.MethodPut, path, nil, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// MassExtendSubscriptionRenewalDate extends the renewal date for all active subscribers
// See https://developer.apple.com/documentation/appstoreserverapi/extend_subscription_renewal_dates_for_all_active_subscribers
func (c *Client) MassExtendSubscriptionRenewalDate(ctx context.Context, req *MassExtendRenewalDateRequest) (*MassExtendRenewalDateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := "/inApps/v1/subscriptions/extend/mass"

	var response MassExtendRenewalDateResponse
	if err := c.makeRequest(ctx, http.MethodPost, path, nil, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMassExtendRenewalDateStatus checks the status of a mass renewal date extension request
// See https://developer.apple.com/documentation/appstoreserverapi/get_status_of_subscription_renewal_date_extensions
func (c *Client) GetMassExtendRenewalDateStatus(ctx context.Context, productID, requestIdentifier string) (*MassExtendRenewalDateStatusResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("productID cannot be empty")
	}
	if requestIdentifier == "" {
		return nil, fmt.Errorf("requestIdentifier cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/mass/%s/%s", productID, requestIdentifier)

	var response MassExtendRenewalDateStatusResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetNotificationHistory gets a list of notifications that the App Store server attempted to send
// See https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
func (c *Client) GetNotificationHistory(ctx context.Context, req *NotificationHistoryRequest) (*NotificationHistoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := "/inApps/v1/notifications/history"

	queryParams := make(url.Values)
	if req.PaginationToken != "" {
		queryParams.Add("paginationToken", req.PaginationToken)
	}

	var response NotificationHistoryResponse
	if err := c.makeRequest(ctx, http.MethodPost, path, queryParams, req, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	for i, v := range response.NotificationHistory {
		payload, err := c.Verifier.VerifyAndDecodeNotification(v.SignedPayload)
		if err != nil {
			return nil, fmt.Errorf("SignedPayload %s\nfailed to verify and decode: %w", v.SignedPayload, err)
		}
		response.NotificationHistory[i].Payload = payload
	}

	return &response, nil
}

// RequestTestNotification asks App Store Server Notifications to send a test notification
// See https://developer.apple.com/documentation/appstoreserverapi/request_a_test_notification
func (c *Client) RequestTestNotification(ctx context.Context) (*SendTestNotificationResponse, error) {
	path := "/inApps/v1/notifications/test"

	var response SendTestNotificationResponse
	if err := c.makeRequest(ctx, http.MethodPost, path, nil, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTestNotificationStatus checks the status of a test notification
// See https://developer.apple.com/documentation/appstoreserverapi/get_test_notification_status
func (c *Client) GetTestNotificationStatus(ctx context.Context, testNotificationToken string) (*CheckTestNotificationResponse, error) {
	if testNotificationToken == "" {
		return nil, fmt.Errorf("testNotificationToken cannot be empty")
	}

	path := fmt.Sprintf("/inApps/v1/notifications/test/%s", testNotificationToken)

	var response CheckTestNotificationResponse
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}

	if c.Verifier.environment == EnvironmentLocalTesting || !c.Verifier.enableAutoDecode {
		return &response, nil
	}

	payload, err := c.Verifier.VerifyAndDecodeNotification(response.SignedPayload)
	if err != nil {
		return nil, fmt.Errorf("SignedPayload %s\nfailed to verify and decode: %w", response.SignedPayload, err)
	}
	response.Payload = payload

	return &response, nil
}

// makereq performs an HTTP req to the App Store Server API
func (c *Client) makeRequest(ctx context.Context, method, path string, queryParams url.Values, requestBody, responseBody any) error {
	token, err := c.TokenGenerator.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate JWT token: %w", err)
	}

	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal req body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP req failed: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return NewAPIErrorFromResponse(resp, respBodyBytes)
	}

	if responseBody != nil {
		if err := json.Unmarshal(respBodyBytes, responseBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
