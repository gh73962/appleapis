package appstoreserver

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gh73962/appleapis/appstoreservernotifications/v2"
)

func TestGetTransactionHistory(t *testing.T) {
	client, err := mockClientWithBody("models/transactionHistoryResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	req := &TransactionHistoryRequest{
		TransactionID:                "1234",
		ProductIDs:                   []string{"com.example.1", "com.example.2"},
		ProductTypes:                 []ProductType{TypeConsumable, TypeAutoRenewableSubscription},
		StartDate:                    time.UnixMilli(1698148800000),
		EndDate:                      time.UnixMilli(1698148900000),
		SubscriptionGroupIdentifiers: []string{"sub_group_id", "sub_group_id_2"},
		InAppOwnershipType:           InAppOwnershipTypeFamilyShared,
		Revoked:                      false,
		Revision:                     "revision_input",
	}
	req.SetSortASC()

	response, err := client.GetTransactionHistory(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.Revision != "revision_output" {
		t.Fatalf("expected %q, got %q", "revision_output", response.Revision)
	}
	if !response.HasMore {
		t.Fatalf("expected %v, got %v", true, response.HasMore)
	}
	if response.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", response.BundleID)
	}
	if response.AppAppleID != 323232 {
		t.Fatalf("expected %v, got %v", 323232, response.AppAppleID)
	}
	if response.Environment != EnvironmentLocalTesting {
		t.Fatalf("expected %q, got %q", EnvironmentLocalTesting, response.Environment)
	}
	if !reflect.DeepEqual(response.SignedTransactions, []string{"signed_transaction_value", "signed_transaction_value2"}) {
		t.Fatalf("expected %v, got %v", []string{"signed_transaction_value", "signed_transaction_value2"}, response.SignedTransactions)
	}
}

func TestGetTransactionInfo(t *testing.T) {
	client, err := mockClientWithBody("models/transactionInfoResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetTransactionInfo(context.Background(), "1234")
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

func TestGetAllSubscriptionStatuses(t *testing.T) {
	client, err := mockClientWithBody("models/getAllSubscriptionStatusesResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetAllSubscriptionStatuses(context.Background(), "4321", StatusActive, StatusExpired)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.Environment != EnvironmentLocalTesting {
		t.Fatalf("expected %q, got %q", EnvironmentLocalTesting, response.Environment)
	}
	if response.BundleID != "com.example" {
		t.Fatalf("expected %q, got %v", "com.example", response.BundleID)
	}
	if response.AppAppleID != 5454545 {
		t.Fatalf("expected %v, got %v", 5454545, response.AppAppleID)
	}
}

func TestSetAppAccountToken(t *testing.T) {
	client, err := mockClientWithBody("", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	req := UpdateAppAccountTokenRequest{
		OriginalTransactionID: "49571273",
		AppAccountToken:       "7389a31a-fb6d-4569-a2a6-db7d85d84813",
	}
	if err := client.SetAppAccountToken(context.Background(), &req); err != nil {
		t.Fatal(err)
	}
}

func TestSetAppAccountTokenError(t *testing.T) {
	client, err := mockClientWithBody("models/invalidAppAccountTokenUUIDError.json", http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}
	req := UpdateAppAccountTokenRequest{
		OriginalTransactionID: "49571273",
		AppAccountToken:       "invalid-uuid",
	}
	if err := client.SetAppAccountToken(context.Background(), &req); err == nil {
		t.Fatal("expected error but got nil")
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.ErrorCode != 4000183 {
			t.Fatalf("expected %v, got %v", 4000183, apiErr.ErrorCode)
		}
	}
}

func TestSetAppAccountTokenNotSupported(t *testing.T) {
	client, err := mockClientWithBody("models/familyTransactionNotSupportedError.json", http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}
	req := UpdateAppAccountTokenRequest{
		OriginalTransactionID: "1234",
		AppAccountToken:       "uuid",
	}
	if err := client.SetAppAccountToken(context.Background(), &req); err == nil {
		t.Fatal("expected error but got nil")
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.ErrorCode != 4000185 {
			t.Fatalf("expected %v, got %v", 4000185, apiErr.ErrorCode)
		}
	}
}

func TestSetAppAccountTokenError2(t *testing.T) {
	client, err := mockClientWithBody("models/transactionIdNotOriginalTransactionId.json", http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}
	req := UpdateAppAccountTokenRequest{
		OriginalTransactionID: "1234",
		AppAccountToken:       "uuid",
	}
	if err := client.SetAppAccountToken(context.Background(), &req); err == nil {
		t.Fatal("expected error but got nil")
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.ErrorCode != 4000187 {
			t.Fatalf("expected %v, got %v", 4000187, apiErr.ErrorCode)
		}
	}
}

func TestSendConsumptionInfo(t *testing.T) {
	client, err := mockClientWithBody("", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	req := ConsumptionRequest{
		TransactionID:            "49571273",
		AccountTenure:            AccountTenure90Days,
		AppAccountToken:          "7389a31a-fb6d-4569-a2a6-db7d85d84813",
		ConsumptionStatus:        ConsumptionStatusNotConsumed,
		CustomerConsented:        true,
		DeliveryStatus:           DeliveryStatusServerOutage,
		LifetimeDollarsPurchased: LifetimeDollarsPurchasedUpTo2000,
		LifetimeDollarsRefunded:  LifetimeDollarsRefundedOver2000,
		Platform:                 PlatformNonApple,
		PlayTime:                 PlayTime4Days,
		RefundPreference:         RefundPreferenceNoPreference,
		SampleContentProvided:    false,
		UserStatus:               UserStatusLimitedAccess,
	}

	if err := client.SendConsumptionInfo(context.Background(), &req); err != nil {
		t.Fatal(err)
	}
}

func TestLookUpOrderID(t *testing.T) {
	client, err := mockClientWithBody("models/lookupOrderIdResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.LookUpOrderID(context.Background(), "W002182")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.Status != OrderLookupStatusInvalid {
		t.Fatalf("expected %q, got %q", OrderLookupStatusInvalid, response.Status)
	}
	if !reflect.DeepEqual(response.SignedTransactions, []string{"signed_transaction_one", "signed_transaction_two"}) {
		t.Fatalf("expected %v, got %v", []string{"signed_transaction_one", "signed_transaction_two"}, response.SignedTransactions)
	}
}

func TestGetRefundHistory(t *testing.T) {
	client, err := mockClientWithBody("models/getRefundHistoryResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetRefundHistory(context.Background(), "555555", "revision_input")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.Revision != "revision_output" {
		t.Fatalf("expected %q, got %q", "revision_output", response.Revision)
	}
	if !response.HasMore {
		t.Fatalf("expected %v, got %v", true, response.HasMore)
	}
	if !reflect.DeepEqual(response.SignedTransactions, []string{"signed_transaction_one", "signed_transaction_two"}) {
		t.Fatalf("expected %v, got %v", []string{"signed_transaction_one", "signed_transaction_two"}, response.SignedTransactions)
	}
}

func TestExtendRenewalDate(t *testing.T) {
	client, err := mockClientWithBody("models/extendSubscriptionRenewalDateResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	request := ExtendRenewalDateRequest{
		OriginalTransactionID: "4124214",
		ExtendByDays:          45,
		ExtendReasonCode:      ExtendReasonCodeCustomerSatisfy,
		RequestIdentifier:     "fdf964a4-233b-486c-aac1-97d8d52688ac",
	}

	response, err := client.ExtendSubscriptionRenewalDate(context.Background(), &request)
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
	client, err := mockClientWithBody("models/extendRenewalDateForAllActiveSubscribersResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	request := &MassExtendRenewalDateRequest{
		ExtendByDays:           45,
		ExtendReasonCode:       ExtendReasonCodeCustomerSatisfy,
		RequestIdentifier:      "fdf964a4-233b-486c-aac1-97d8d52688ac",
		ProductID:              "com.example.productId",
		StorefrontCountryCodes: []string{"USA", "MEX"},
	}

	response, err := client.MassExtendSubscriptionRenewalDate(context.Background(), request)
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

func TestGetMassExtendRenewalDateStatus(t *testing.T) {
	client, err := mockClientWithBody("models/getStatusOfSubscriptionRenewalDateExtensionsResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetMassExtendRenewalDateStatus(context.Background(), "com.example.product", "20fba8a0-2b80-4a7d-a17f-85c1854727f8")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.RequestIdentifier != "20fba8a0-2b80-4a7d-a17f-85c1854727f8" {
		t.Fatalf("expected %q, got %q", "20fba8a0-2b80-4a7d-a17f-85c1854727f8", response.RequestIdentifier)
	}
	if !response.Complete {
		t.Fatal("expected true")
	}
	if response.CompleteDate != 1698148900000 {
		t.Fatalf("expected %v, got %v", 1698148900000, response.CompleteDate)
	}
	if response.SucceededCount != 30 {
		t.Fatalf("expected %v, got %v", 30, response.SucceededCount)
	}
	if response.FailedCount != 2 {
		t.Fatalf("expected %v, got %v", 2, response.FailedCount)
	}
}

func TestGetNotificationHistory(t *testing.T) {
	client, err := mockClientWithBody("models/getNotificationHistoryResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	request := &NotificationHistoryRequest{
		StartTime:           time.UnixMilli(1698148900000),
		EndTime:             time.UnixMilli(1698148950000),
		NotificationType:    appstoreservernotifications.TypeSubscribed,
		NotificationSubtype: appstoreservernotifications.SubtypeInitialBuy,
		OnlyFailures:        true,
		TransactionID:       "999733843",
	}

	response, err := client.GetNotificationHistory(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.PaginationToken != "57715481-805a-4283-8499-1c19b5d6b20a" {
		t.Fatalf("expected %q, got %q", "57715481-805a-4283-8499-1c19b5d6b20a", response.PaginationToken)
	}
	if !response.HasMore {
		t.Fatal("expected true")
	}

	cmpData := []NotificationHistoryResponseItem{
		{
			SignedPayload: "signed_payload_one",
			SendAttempts: []SendAttemptItem{
				{
					AttemptDate:       1698148900000,
					SendAttemptResult: SendAttemptResultNoResponse,
				},
				{
					AttemptDate:       1698148950000,
					SendAttemptResult: SendAttemptResultSuccess,
				},
			},
		},
		{
			SignedPayload: "signed_payload_two",
			SendAttempts: []SendAttemptItem{
				{
					AttemptDate:       1698148800000,
					SendAttemptResult: SendAttemptResultCircularRedirect,
				},
			},
		},
	}
	if !reflect.DeepEqual(response.NotificationHistory, cmpData) {
		t.Fatalf("expected %v, got %v", cmpData, response.NotificationHistory)
	}
}

func TestRequestTestNotification(t *testing.T) {
	client, err := mockClientWithBody("models/requestTestNotificationResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.RequestTestNotification(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.TestNotificationToken != "ce3af791-365e-4c60-841b-1674b43c1609" {
		t.Fatalf("expected %q, got %q", "ce3af791-365e-4c60-841b-1674b43c1609", response.TestNotificationToken)
	}
}

func TestGetTestNotificationStatus(t *testing.T) {
	client, err := mockClientWithBody("models/getTestNotificationStatusResponse.json", http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetTestNotificationStatus(context.Background(), "8cd2974c-f905-4da0-bf52-7f12e55184be")
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("expected non-nil response")
	}
	if response.SignedPayload != "signed_payload" {
		t.Fatal("expected non-empty SignedPayload")
	}
	cmpData := []SendAttemptItem{
		{
			AttemptDate:       1698148900000,
			SendAttemptResult: SendAttemptResultNoResponse,
		},
		{
			AttemptDate:       1698148950000,
			SendAttemptResult: SendAttemptResultSuccess,
		},
	}
	if !reflect.DeepEqual(response.SendAttempts, cmpData) {
		t.Fatalf("expected %v, got %v", cmpData, response.SendAttempts)
	}
}

func TestAPIErrorHandling(t *testing.T) {
	client, err := mockClientWithBody("models/apiException.json", http.StatusInternalServerError)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RequestTestNotification(context.Background())
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected error to be of type *APIError but got %T", err)
	}
	if apiErr.ErrorCode != 5000000 {
		t.Fatalf("expected error code %d but got %d", 5000000, apiErr.ErrorCode)
	}
}

func TestAPITooManyRequestsException(t *testing.T) {
	client, err := mockClientWithBody("models/apiTooManyRequestsException.json", http.StatusTooManyRequests)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RequestTestNotification(context.Background())
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected error to be of type *APIError but got %T", err)
	}
	if apiErr.ErrorCode != 4290000 {
		t.Fatalf("expected error code %d but got %d", 4290000, apiErr.ErrorCode)
	}
}
