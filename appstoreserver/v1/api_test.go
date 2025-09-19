package appstoreserver

import (
	"context"
	"strings"
	"testing"
)

func TestGetTransactionInfo(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("GET", "/inApps/v1/transactions/1233214", "models/transactionInfoResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Test the API call
	response, err := client.client.GetTransactionInfo(context.Background(), "1233214")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.SignedTransactionInfo != "signed_transaction_info_value" {
		t.Fatalf("expected %q, got %q", "signed_transaction_info_value", response.SignedTransactionInfo)
	}
}

func TestGetTransactionHistory(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("GET", "/inApps/v2/history/1233214", "models/transactionHistoryResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Create request
	request := &TransactionHistoryRequest{
		Sort:         "ASCENDING",
		ProductIDs:   []string{"com.example.1", "com.example.2"},
		ProductTypes: []string{"Auto-Renewable Subscription"},
		StartDate:    1698148800000,
		EndDate:      1698148900000,
	}

	// Test the API call
	response, err := client.client.GetTransactionHistory(context.Background(), "1233214", request)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", response.BundleID)
	}
	if response.AppAppleID != 1234 {
		t.Fatalf("expected %v, got %v", 1234, response.AppAppleID)
	}
	if response.SignedTransactions == nil {
		t.Fatal("expected non-nil SignedTransactions")
	}
	if len(response.SignedTransactions) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(response.SignedTransactions))
	}
}

func TestLookUpOrderID(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("GET", "/inApps/v1/lookup/W002182", "models/lookupOrderIdResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Test the API call
	response, err := client.client.LookUpOrderID(context.Background(), "W002182")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.SignedTransactions == nil {
		t.Fatal("expected non-nil SignedTransactions")
	}
	if len(response.SignedTransactions) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(response.SignedTransactions))
	}
}

func TestGetRefundHistory(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("GET", "/inApps/v2/refund/lookup/555555", "models/getRefundHistoryResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Test the API call
	response, err := client.client.GetRefundHistory(context.Background(), "555555", nil)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.SignedTransactions == nil {
		t.Fatal("expected non-nil SignedTransactions")
	}
	if len(response.SignedTransactions) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(response.SignedTransactions))
	}
	if response.SignedTransactions[0] != "signed_transaction_one" {
		t.Fatalf("expected %q, got %q", "signed_transaction_one", response.SignedTransactions[0])
	}
}

func TestExtendRenewalDate(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("PUT", "/inApps/v1/subscriptions/extend/4124214", "models/extendSubscriptionRenewalDateResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Create request
	request := &ExtendRenewalDateRequest{
		ExtendByDays:      45,
		ExtendReasonCode:  1, // CUSTOMER_SATISFACTION
		RequestIdentifier: "fdf964a4-233b-486c-aac1-97d8d52688ac",
	}

	// Test the API call
	response, err := client.client.ExtendRenewalDate(context.Background(), "4124214", request)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.OriginalTransactionID != "2312412" {
		t.Fatalf("expected %q, got %q", "2312412", response.OriginalTransactionID)
	}
	if response.WebOrderLineItemID != "9993" {
		t.Fatalf("expected %q, got %q", "9993", response.WebOrderLineItemID)
	}
	if !response.Success {
		t.Fatal("expected true")
	}
	if response.EffectiveDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), response.EffectiveDate)
	}
}

func TestMassExtendRenewalDate(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("POST", "/inApps/v1/subscriptions/extend/mass", "models/extendRenewalDateForAllActiveSubscribersResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Create request
	request := &MassExtendRenewalDateRequest{
		ExtendByDays:      45,
		ExtendReasonCode:  1, // CUSTOMER_SATISFACTION
		RequestIdentifier: "fdf964a4-233b-486c-aac1-97d8d52688ac",
		ProductID:         "com.example.productId",
	}

	// Test the API call
	response, err := client.client.MassExtendRenewalDate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.RequestIdentifier != "758883e8-151b-47b7-abd0-60c4d804c2f5" {
		t.Fatalf("expected %q, got %q", "758883e8-151b-47b7-abd0-60c4d804c2f5", response.RequestIdentifier)
	}
}

func TestRequestTestNotification(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("POST", "/inApps/v1/notifications/test", "models/requestTestNotificationResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Test the API call
	response, err := client.client.RequestTestNotification(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.TestNotificationToken != "ce3af791-365b-4c60-841b-1674b43c1609" {
		t.Fatalf("expected %q, got %q", "ce3af791-365b-4c60-841b-1674b43c1609", response.TestNotificationToken)
	}
}

func TestGetTestNotificationStatus(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("GET", "/inApps/v1/notifications/test/8cd2974c-f905-4da0-bf52-7f12e55184be", "models/getTestNotificationStatusResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Test the API call
	response, err := client.client.GetTestNotificationStatus(context.Background(), "8cd2974c-f905-4da0-bf52-7f12e55184be")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.SignedPayload == "" {
		t.Fatal("expected non-empty SignedPayload")
	}
}

func TestGetNotificationHistory(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up mock response
	err := helper.SetResponseFromFile("POST", "/inApps/v1/notifications/history", "models/getNotificationHistoryResponse.json")
	if err != nil {
		t.Fatal(err)
	}

	client := helper.CreateTestClient(t)

	// Create request
	notificationType := "TEST"
	notificationSubtype := "INITIAL_BUY"
	onlyFailures := false
	transactionID := "1233214"

	request := &NotificationHistoryRequest{
		StartDate:           1698148800000,
		EndDate:             1698148900000,
		NotificationType:    &notificationType,
		NotificationSubtype: &notificationSubtype,
		OnlyFailures:        &onlyFailures,
		TransactionID:       &transactionID,
	}

	// Test the API call
	response, err := client.client.GetNotificationHistory(context.Background(), nil, request)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.NotificationHistory == nil {
		t.Fatal("expected non-nil NotificationHistory")
	}
	if len(response.NotificationHistory) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(response.NotificationHistory))
	}
	if response.NotificationHistory[0].SignedPayload != "signed_payload" {
		t.Fatalf("expected %q, got %q", "signed_payload", response.NotificationHistory[0].SignedPayload)
	}
}

func TestAPIErrorHandling(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up error response
	helper.SetResponse("GET", "/inApps/v1/transactions/invalid", TestResponse{
		StatusCode: 400,
		Body:       `{"errorCode": 4040001, "errorMessage": "An error occurred"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	})

	client := helper.CreateTestClient(t)

	// Test error handling
	_, err := client.client.GetTransactionInfo(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "4040001") {
		t.Fatalf("expected error to contain %q but got %q", "4040001", err.Error())
	}
}

func TestAPITooManyRequestsException(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Set up rate limit response
	helper.SetResponse("GET", "/inApps/v1/transactions/ratelimited", TestResponse{
		StatusCode: 429,
		Body:       `{"errorCode": 4290000, "errorMessage": "Rate limit exceeded"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	})

	client := helper.CreateTestClient(t)

	// Test rate limit handling
	_, err := client.client.GetTransactionInfo(context.Background(), "ratelimited")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "4290000") {
		t.Fatalf("expected error to contain %q but got %q", "4290000", err.Error())
	}
}
