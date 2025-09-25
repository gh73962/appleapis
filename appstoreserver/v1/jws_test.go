package appstoreserver

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gh73962/appleapis/appstoreservernotifications/v2"
)

func TestTransactionDecoding(t *testing.T) {
	client, err := mockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	signedTransaction, err := mockSignedData("models/signedTransaction.json")
	if err != nil {
		t.Fatal(err)
	}

	decodedPayload, err := client.verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	if decodedPayload.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.BundleID)
	}
	if decodedPayload.ProductID != "com.example.product" {
		t.Fatalf("expected %q, got %q", "com.example.product", decodedPayload.ProductID)
	}
	if decodedPayload.SubscriptionGroupIdentifier != nil {
		if *decodedPayload.SubscriptionGroupIdentifier != "55555" {
			t.Fatalf("expected %q, got %q", "55555", *decodedPayload.SubscriptionGroupIdentifier)
		}
	}
	if decodedPayload.PurchaseDate != 1698148900000 {
		t.Fatalf("expected %v, got %v", 1698148900000, decodedPayload.PurchaseDate)
	}
	if decodedPayload.OriginalPurchaseDate != 1698148800000 {
		t.Fatalf("expected %v, got %v", 1698148800000, decodedPayload.OriginalPurchaseDate)
	}
	if decodedPayload.ExpiresDate != nil {
		if *decodedPayload.ExpiresDate != int64(1698149000000) {
			t.Fatalf("expected %v, got %v", int64(1698149000000), *decodedPayload.ExpiresDate)
		}
	}
	if decodedPayload.Quantity != 1 {
		t.Fatalf("expected %v, got %v", 1, decodedPayload.Quantity)
	}
	if string(decodedPayload.Type) != string(TypeAutoRenewableSubscription) {
		t.Fatalf("expected %q, got %q", "Auto-Renewable Subscription", string(decodedPayload.Type))
	}
	if decodedPayload.AppAccountToken != nil {
		if *decodedPayload.AppAccountToken != "7e3fb20b-4cdb-47cc-936d-99d65f608138" {
			t.Fatalf("expected %q, got %q", "7e3fb20b-4cdb-47cc-936d-99d65f608138", *decodedPayload.AppAccountToken)
		}
	}
	if string(decodedPayload.InAppOwnershipType) != "PURCHASED" {
		t.Fatalf("expected %q, got %q", "PURCHASED", string(decodedPayload.InAppOwnershipType))
	}
	if decodedPayload.SignedDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), decodedPayload.SignedDate)
	}
	if decodedPayload.RevocationDate != nil {
		if *decodedPayload.RevocationDate != int64(1698148950000) {
			t.Fatalf("expected %v, got %v", int64(1698148950000), *decodedPayload.RevocationDate)
		}
	}
	if decodedPayload.RevocationReason != nil {
		if *decodedPayload.RevocationReason != RevocationReasonAppIssue {
			t.Fatalf("expected %v, got %v", RevocationReasonAppIssue, *decodedPayload.RevocationReason)
		}
	}
	if decodedPayload.IsUpgraded != nil {
		if !*decodedPayload.IsUpgraded {
			t.Fatal("expected true")
		}
	}
	if decodedPayload.OfferType != nil {
		if *decodedPayload.OfferType != OfferTypeIntroductory {
			t.Fatalf("expected %v, got %v", OfferTypeIntroductory, *decodedPayload.OfferType)
		}
	}
	if decodedPayload.OfferIdentifier != nil {
		if *decodedPayload.OfferIdentifier != "abc.123" {
			t.Fatalf("expected %q, got %q", "abc.123", *decodedPayload.OfferIdentifier)
		}
	}
	if decodedPayload.TransactionReason != nil {
		if string(*decodedPayload.TransactionReason) != "PURCHASE" {
			t.Fatalf("expected %q, got %q", "PURCHASE", string(*decodedPayload.TransactionReason))
		}
	}
}

func TestRenewalInfoDecoding(t *testing.T) {
	client, err := mockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	signedRenewalInfo, err := mockSignedData("models/signedRenewalInfo.json")
	if err != nil {
		t.Fatal(err)
	}

	decodedPayload, err := client.verifier.VerifyAndDecodeRenewalInfo(signedRenewalInfo)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	if decodedPayload.OriginalTransactionID != "12345" {
		t.Fatalf("expected %q, got %q", "12345", decodedPayload.OriginalTransactionID)
	}
	if decodedPayload.ProductID != "com.example.product" {
		t.Fatalf("expected %q, got %q", "com.example.product", decodedPayload.ProductID)
	}
	if decodedPayload.AutoRenewStatus != AutoRenewStatusOn {
		t.Fatalf("expected %v, got %v", AutoRenewStatusOn, decodedPayload.AutoRenewStatus)
	}
	if decodedPayload.ExpirationIntent != nil {
		if *decodedPayload.ExpirationIntent != ExpirationIntentCustomerCanceled {
			t.Fatalf("expected %v, got %v", ExpirationIntentCustomerCanceled, *decodedPayload.ExpirationIntent)
		}
	}
	if decodedPayload.SignedDate != int64(1698148800000) {
		t.Fatalf("expected %v, got %v", int64(1698148800000), decodedPayload.SignedDate)
	}
	if string(decodedPayload.Environment) != "LocalTesting" {
		t.Fatalf("expected %q, got %q", "LocalTesting", string(decodedPayload.Environment))
	}
	if decodedPayload.RecentSubscriptionStartDate != nil {
		if *decodedPayload.RecentSubscriptionStartDate != int64(1698148800000) {
			t.Fatalf("expected %v, got %v", int64(1698148800000), *decodedPayload.RecentSubscriptionStartDate)
		}
	}
	if decodedPayload.OfferType != nil {
		if *decodedPayload.OfferType != OfferTypePromotional {
			t.Fatalf("expected %v, got %v", OfferTypePromotional, *decodedPayload.OfferType)
		}
	}
	if decodedPayload.OfferIdentifier != nil {
		if *decodedPayload.OfferIdentifier != "abc.123" {
			t.Fatalf("expected %q, got %q", "abc.123", *decodedPayload.OfferIdentifier)
		}
	}
}

func TestNotificationDecoding(t *testing.T) {
	client, err := mockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	signedNotification, err := mockSignedData("models/signedNotification.json")
	if err != nil {
		t.Fatal(err)
	}

	decodedPayload, err := client.verifier.VerifyAndDecodeNotification(signedNotification)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	if string(decodedPayload.NotificationType) != "SUBSCRIBED" {
		t.Fatalf("expected %q, got %q", "SUBSCRIBED", string(decodedPayload.NotificationType))
	}
	if decodedPayload.Subtype != "INITIAL_BUY" {
		t.Fatalf("expected %q, got %q", "INITIAL_BUY", decodedPayload.Subtype)
	}
	if decodedPayload.NotificationUUID != "002e14d5-51f5-4503-b5a8-c3a1af68eb20" {
		t.Fatalf("expected %q, got %q", "002e14d5-51f5-4503-b5a8-c3a1af68eb20", decodedPayload.NotificationUUID)
	}
	if decodedPayload.Version != "2.0" {
		t.Fatalf("expected %q, got %q", "2.0", decodedPayload.Version)
	}
	if decodedPayload.SignedDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), decodedPayload.SignedDate)
	}
	if decodedPayload.Data == nil {
		t.Fatal("expected non-nil Data")
	}
	if string(decodedPayload.Data.Environment) != "LocalTesting" {
		t.Fatalf("expected %q, got %q", "LocalTesting", string(decodedPayload.Data.Environment))
	}
	if decodedPayload.Data.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.Data.BundleID)
	}
	if decodedPayload.Data.AppAppleID != int64(41234) {
		t.Fatalf("expected %v, got %v", int64(41234), decodedPayload.Data.AppAppleID)
	}
	if decodedPayload.Data.SignedTransactionInfo != "signed_transaction_info_value" {
		t.Fatalf("expected %q, got %q", "signed_transaction_info_value", decodedPayload.Data.SignedTransactionInfo)
	}
	if decodedPayload.Data.SignedRenewalInfo != "signed_renewal_info_value" {
		t.Fatalf("expected %q, got %q", "signed_renewal_info_value", decodedPayload.Data.SignedRenewalInfo)
	}
}

func TestConsumptionRequestNotificationDecoding(t *testing.T) {
	client, err := mockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	signedNotification, err := mockSignedData("models/signedConsumptionRequestNotification.json")
	if err != nil {
		t.Fatal(err)
	}

	decodedPayload, err := client.verifier.VerifyAndDecodeNotification(signedNotification)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	// Verify the decoded fields
	if string(decodedPayload.NotificationType) != "CONSUMPTION_REQUEST" {
		t.Fatalf("expected %q, got %q", "CONSUMPTION_REQUEST", string(decodedPayload.NotificationType))
	}
	if decodedPayload.Subtype != "" {
		t.Fatalf("expected empty subtype, got %q", decodedPayload.Subtype)
	}
	if decodedPayload.NotificationUUID != "002e14d5-51f5-4503-b5a8-c3a1af68eb20" {
		t.Fatalf("expected %q, got %q", "002e14d5-51f5-4503-b5a8-c3a1af68eb20", decodedPayload.NotificationUUID)
	}
	if decodedPayload.Version != "2.0" {
		t.Fatalf("expected %q, got %q", "2.0", decodedPayload.Version)
	}
	if decodedPayload.SignedDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), decodedPayload.SignedDate)
	}
	if decodedPayload.Data == nil {
		t.Fatal("expected non-nil Data")
	}
	if decodedPayload.Summary != nil {
		t.Fatal("expected nil Summary")
	}
	if decodedPayload.ExternalPurchaseToken != nil {
		t.Fatal("expected nil ExternalPurchaseToken")
	}
	if string(decodedPayload.Data.Environment) != "LocalTesting" {
		t.Fatalf("expected %q, got %q", "LocalTesting", string(decodedPayload.Data.Environment))
	}
	if decodedPayload.Data.AppAppleID != int64(41234) {
		t.Fatalf("expected %v, got %v", int64(41234), decodedPayload.Data.AppAppleID)
	}
	if decodedPayload.Data.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.Data.BundleID)
	}
	if decodedPayload.Data.SignedTransactionInfo != "signed_transaction_info_value" {
		t.Fatalf("expected %q, got %q", "signed_transaction_info_value", decodedPayload.Data.SignedTransactionInfo)
	}
	if decodedPayload.Data.SignedRenewalInfo != "signed_renewal_info_value" {
		t.Fatalf("expected %q, got %q", "signed_renewal_info_value", decodedPayload.Data.SignedRenewalInfo)
	}

}

func TestSummaryNotificationDecoding(t *testing.T) {
	client, err := mockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	signedNotification, err := mockSignedData("models/signedSummaryNotification.json")
	if err != nil {
		t.Fatal(err)
	}

	decodedPayload, err := client.verifier.VerifyAndDecodeNotification(signedNotification)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	if decodedPayload.NotificationType != appstoreservernotifications.TypeRenewalExtension {
		t.Fatalf("expected %q, got %q", appstoreservernotifications.TypeRenewalExtension, decodedPayload.NotificationType)
	}
	if decodedPayload.Subtype != appstoreservernotifications.SubtypeSummary {
		t.Fatalf("expected %q, got %q", appstoreservernotifications.SubtypeSummary, decodedPayload.Subtype)
	}
	if decodedPayload.NotificationUUID != "002e14d5-51f5-4503-b5a8-c3a1af68eb20" {
		t.Fatalf("expected %q, got %q", "002e14d5-51f5-4503-b5a8-c3a1af68eb20", decodedPayload.NotificationUUID)
	}
	if decodedPayload.Version != "2.0" {
		t.Fatalf("expected %q, got %q", "2.0", decodedPayload.Version)
	}
	if decodedPayload.SignedDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), decodedPayload.SignedDate)
	}
	if decodedPayload.Data != nil {
		t.Fatal("expected nil Data")
	}
	if decodedPayload.Summary == nil {
		t.Fatal("expected non-nil Summary")
	}
	if decodedPayload.ExternalPurchaseToken != nil {
		t.Fatal("expected nil ExternalPurchaseToken")
	}
	if string(decodedPayload.Summary.Environment) != "LocalTesting" {
		t.Fatalf("expected %q, got %q", "LocalTesting", string(decodedPayload.Summary.Environment))
	}
	if decodedPayload.Summary.AppAppleID != int64(41234) {
		t.Fatalf("expected %v, got %v", int64(41234), decodedPayload.Summary.AppAppleID)
	}
	if decodedPayload.Summary.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.Summary.BundleID)
	}
	if decodedPayload.Summary.ProductID != "com.example.product" {
		t.Fatalf("expected %q, got %q", "com.example.product", decodedPayload.Summary.ProductID)
	}
	if decodedPayload.Summary.RequestIdentifier != "efb27071-45a4-4aca-9854-2a1e9146f265" {
		t.Fatalf("expected %q, got %q", "efb27071-45a4-4aca-9854-2a1e9146f265", decodedPayload.Summary.RequestIdentifier)
	}
	if len(decodedPayload.Summary.StorefrontCountryCodes) != 3 {
		t.Fatalf("expected 3 country codes, got %d", len(decodedPayload.Summary.StorefrontCountryCodes))
	}
	expectedCountries := []string{"CAN", "USA", "MEX"}
	if !reflect.DeepEqual(decodedPayload.Summary.StorefrontCountryCodes, expectedCountries) {
		t.Fatalf("expected %v, got %v", expectedCountries, decodedPayload.Summary.StorefrontCountryCodes)
	}
	if decodedPayload.Summary.SucceededCount != 5 {
		t.Fatalf("expected %v, got %v", 5, decodedPayload.Summary.SucceededCount)
	}
	if decodedPayload.Summary.FailedCount != 2 {
		t.Fatalf("expected %v, got %v", 2, decodedPayload.Summary.FailedCount)
	}
}

func TestExternalPurchaseTokenNotificationDecoding(t *testing.T) {
	t.Skip("TODO")
}

func TestExternalPurchaseTokenSandboxNotificationDecoding(t *testing.T) {
	t.Skip("TODO")
}

// class PayloadVerification(unittest.TestCase)
func TestAppStoreServerNotificationDecoding(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentSandbox))
	if err != nil {
		t.Fatal(err)
	}
	testNotification, err := os.ReadFile("../../testdata/mock_signed_data/testNotification")
	if err != nil {
		t.Fatal(err)
	}

	// Test notification decoding
	decodedPayload, err := client.verifier.VerifyAndDecodeNotification(string(testNotification))
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil value")
	}

	if decodedPayload.NotificationType != appstoreservernotifications.TypeTest {
		t.Fatalf("expected %q, got %q", appstoreservernotifications.TypeTest, decodedPayload.NotificationType)
	}
}

func TestAppStoreServerNotificationDecodingProduction(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentProduction))
	if err != nil {
		t.Fatal(err)
	}
	testNotification, err := os.ReadFile("../../testdata/mock_signed_data/testNotification")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.verifier.VerifyAndDecodeNotification(string(testNotification))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestMissingX5CHeader(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentSandbox))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("../../testdata/mock_signed_data/missingX5CHeaderClaim")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.verifier.VerifyAndDecodeNotification(string(data))
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if !strings.Contains(err.Error(), "x5c") {
		t.Fatalf("Received expected error for missing x5c header: %v", err)
	}
}

func TestWrongBundleIDForServerNotification(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentSandbox), WithBundleID("com.examplex"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("../../testdata/mock_signed_data/wrongBundleId")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.verifier.VerifyAndDecodeNotification(string(data))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "INVALID_APP_IDENTIFIER") {
		t.Fatalf("expected error to contain %q but got %q", "bundle", err.Error())
	}
}

func TestWrongAppAppleIDForServerNotification(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentProduction), WithAppAppleID(1235))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("../../testdata/mock_signed_data/testNotification")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.verifier.VerifyAndDecodeNotification(string(data))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "INVALID_APP_IDENTIFIER") {
		t.Fatalf("expected error to contain %q but got %q", "INVALID_APP_IDENTIFIER", err.Error())
	}
}

func TestMalformedJWTWithTooManyParts(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentSandbox))
	if err != nil {
		t.Fatal(err)
	}

	malformedJWT := "header.payload.signature.extra"
	_, err = client.verifier.VerifyAndDecodeSignedTransaction(malformedJWT)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if !strings.Contains(err.Error(), "token is malformed") {
		t.Fatalf("Received expected error for invalid JWT: %v", err)
	}
}

func TestMalformedJWTWithMalformedData(t *testing.T) {
	client, err := mockTestClient(WithEnvironment(EnvironmentSandbox))
	if err != nil {
		t.Fatal(err)
	}

	malformedJWT := "not.a.jwt"

	_, err = client.verifier.VerifyAndDecodeSignedTransaction(malformedJWT)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "token is malformed") {
		t.Fatalf("Received expected error for invalid JWT: %v", err)
	}
}
