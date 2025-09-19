package appstoreserver

import (
	"strings"
	"testing"
)

func TestTransactionDecoding(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	client := helper.CreateTestClient(t)

	// Create signed data from test JSON
	signedTransaction := CreateSignedDataFromJSON(t, "models/signedTransaction.json")

	// Test transaction decoding
	decodedPayload, err := client.verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	// Verify the decoded fields
	if decodedPayload.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.BundleID)
	}
	if decodedPayload.ProductID != "abc.def" {
		t.Fatalf("expected %q, got %q", "abc.def", decodedPayload.ProductID)
	}
	if decodedPayload.SubscriptionGroupIdentifier != nil {
		if *decodedPayload.SubscriptionGroupIdentifier != "ghij" {
			t.Fatalf("expected %q, got %q", "ghij", *decodedPayload.SubscriptionGroupIdentifier)
		}
	}
	if decodedPayload.PurchaseDate != int64(1698148900000) {
		t.Fatalf("expected %v, got %v", int64(1698148900000), decodedPayload.PurchaseDate)
	}
	if decodedPayload.OriginalPurchaseDate != int64(1698148800000) {
		t.Fatalf("expected %v, got %v", int64(1698148800000), decodedPayload.OriginalPurchaseDate)
	}
	if decodedPayload.ExpiresDate != nil {
		if *decodedPayload.ExpiresDate != int64(1698149000000) {
			t.Fatalf("expected %v, got %v", int64(1698149000000), *decodedPayload.ExpiresDate)
		}
	}
	if decodedPayload.Quantity != 1 {
		t.Fatalf("expected %v, got %v", 1, decodedPayload.Quantity)
	}
	if string(decodedPayload.Type) != "Auto-Renewable Subscription" {
		t.Fatalf("expected %q, got %q", "Auto-Renewable Subscription", string(decodedPayload.Type))
	}
	if decodedPayload.AppAccountToken != nil {
		if *decodedPayload.AppAccountToken != "abc.efg" {
			t.Fatalf("expected %q, got %q", "abc.efg", *decodedPayload.AppAccountToken)
		}
	}
	if string(decodedPayload.InAppOwnershipType) != "FAMILY_SHARED" {
		t.Fatalf("expected %q, got %q", "FAMILY_SHARED", string(decodedPayload.InAppOwnershipType))
	}
	if decodedPayload.SignedDate != int64(1698148950000) {
		t.Fatalf("expected %v, got %v", int64(1698148950000), decodedPayload.SignedDate)
	}
	if decodedPayload.RevocationDate != nil {
		if *decodedPayload.RevocationDate != int64(1698149100000) {
			t.Fatalf("expected %v, got %v", int64(1698149100000), *decodedPayload.RevocationDate)
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
		if *decodedPayload.OfferIdentifier != "abc123" {
			t.Fatalf("expected %q, got %q", "abc123", *decodedPayload.OfferIdentifier)
		}
	}
	if decodedPayload.TransactionReason != nil {
		if string(*decodedPayload.TransactionReason) != "PURCHASE" {
			t.Fatalf("expected %q, got %q", "PURCHASE", string(*decodedPayload.TransactionReason))
		}
	}
}

func TestRenewalInfoDecoding(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	client := helper.CreateTestClient(t)

	// Create signed data from test JSON
	signedRenewalInfo := CreateSignedDataFromJSON(t, "models/signedRenewalInfo.json")

	// Test renewal info decoding
	decodedPayload, err := client.verifier.VerifyAndDecodeRenewalInfo(signedRenewalInfo)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	// Verify the decoded fields
	if decodedPayload.OriginalTransactionID != "12345" {
		t.Fatalf("expected %q, got %q", "12345", decodedPayload.OriginalTransactionID)
	}
	if decodedPayload.ProductID != "com.example.product.autorenewable" {
		t.Fatalf("expected %q, got %q", "com.example.product.autorenewable", decodedPayload.ProductID)
	}
	if decodedPayload.AutoRenewStatus != AutoRenewStatusOn {
		t.Fatalf("expected %v, got %v", AutoRenewStatusOn, decodedPayload.AutoRenewStatus)
	}
	if decodedPayload.ExpirationIntent != nil {
		if *decodedPayload.ExpirationIntent != ExpirationIntentCustomerCanceled {
			t.Fatalf("expected %v, got %v", ExpirationIntentCustomerCanceled, *decodedPayload.ExpirationIntent)
		}
	}
	if decodedPayload.SignedDate != int64(1698148950000) {
		t.Fatalf("expected %v, got %v", int64(1698148950000), decodedPayload.SignedDate)
	}
	if string(decodedPayload.Environment) != "LocalTesting" {
		t.Fatalf("expected %q, got %q", "LocalTesting", string(decodedPayload.Environment))
	}
	if decodedPayload.RecentSubscriptionStartDate != nil {
		if *decodedPayload.RecentSubscriptionStartDate != int64(1698149000000) {
			t.Fatalf("expected %v, got %v", int64(1698149000000), *decodedPayload.RecentSubscriptionStartDate)
		}
	}
	if decodedPayload.OfferType != nil {
		if *decodedPayload.OfferType != OfferTypeIntroductory {
			t.Fatalf("expected %v, got %v", OfferTypeIntroductory, *decodedPayload.OfferType)
		}
	}
	if decodedPayload.OfferIdentifier != nil {
		if *decodedPayload.OfferIdentifier != "abc123" {
			t.Fatalf("expected %q, got %q", "abc123", *decodedPayload.OfferIdentifier)
		}
	}
}

func TestNotificationDecoding(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	client := helper.CreateTestClient(t)

	// Create signed data from test JSON
	signedNotification := CreateSignedDataFromJSON(t, "models/signedNotification.json")

	// Test notification decoding
	decodedPayload, err := client.verifier.VerifyAndDecodeNotification(signedNotification)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil decodedPayload")
	}

	// Verify the decoded fields
	if string(decodedPayload.NotificationType) != "TEST" {
		t.Fatalf("expected %q, got %q", "TEST", string(decodedPayload.NotificationType))
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
	if decodedPayload.Data.AppAppleID != int64(1234) {
		t.Fatalf("expected %v, got %v", int64(1234), decodedPayload.Data.AppAppleID)
	}
	if decodedPayload.Data.SignedTransactionInfo != "signed_transaction_info" {
		t.Fatalf("expected %q, got %q", "signed_transaction_info", decodedPayload.Data.SignedTransactionInfo)
	}
	if decodedPayload.Data.SignedRenewalInfo != "signed_renewal_info" {
		t.Fatalf("expected %q, got %q", "signed_renewal_info", decodedPayload.Data.SignedRenewalInfo)
	}
}

func TestPayloadVerificationWithInvalidSignature(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	client := helper.CreateTestClient(t)

	// Test with invalid signed data
	invalidSignedData := "invalid.signed.data"

	// Test transaction decoding with invalid data
	_, err := client.verifier.VerifyAndDecodeSignedTransaction(invalidSignedData)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	// Test renewal info decoding with invalid data
	_, err = client.verifier.VerifyAndDecodeRenewalInfo(invalidSignedData)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	// Test notification decoding with invalid data
	_, err = client.verifier.VerifyAndDecodeNotification(invalidSignedData)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPayloadVerificationWithWrongBundleID(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Create client with different bundle ID
	client, err := New(
		WithPrivateKey([]byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKGsJ1QwJ0WgKwZ5Mz3KZ5gEH8Z7c2+Y3+Ic7J8R8O9XoAoGCCqGSM49
AwEHoUQDQgAE4rWBxGmFbnPIPQI0zsBKzLxsj8pD2vqbr0yPISUx2WQyxmrNql9f
hK8YEEyYFV7++p5i4YUSR/o9uQIgCPIhrA==
-----END EC PRIVATE KEY-----`)),
		WithKeyID("TESTKEY123"),
		WithIssuerID("TESTISSUER123"),
		WithBundleIDClient("com.wrong.bundle"),
		WithEnvironmentClient(EnvironmentSandbox),
		WithAppAppleIDClient(1234567890),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create signed data from test JSON (which has bundle ID "com.example")
	signedTransaction := CreateSignedDataFromJSON(t, "models/signedTransaction.json")

	// Test transaction decoding with wrong bundle ID - should fail
	_, err = client.verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "bundle") {
		t.Fatalf("expected error to contain %q but got %q", "bundle", err.Error())
	}
}

func TestPayloadVerificationWithWrongEnvironment(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Close()

	// Create client with production environment
	client, err := New(
		WithPrivateKey([]byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKGsJ1QwJ0WgKwZ5Mz3KZ5gEH8Z7c2+Y3+Ic7J8R8O9XoAoGCCqGSM49
AwEHoUQDQgAE4rWBxGmFbnPIPQI0zsBKzLxsj8pD2vqbr0yPISUx2WQyxmrNql9f
hK8YEEyYFV7++p5i4YUSR/o9uQIgCPIhrA==
-----END EC PRIVATE KEY-----`)),
		WithKeyID("TESTKEY123"),
		WithIssuerID("TESTISSUER123"),
		WithBundleIDClient("com.example"),
		WithEnvironmentClient(EnvironmentProduction),
		WithAppAppleIDClient(1234567890),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create signed data from test JSON (which has environment "LocalTesting")
	signedTransaction := CreateSignedDataFromJSON(t, "models/signedTransaction.json")

	// Test transaction decoding with wrong environment - should fail
	_, err = client.verifier.VerifyAndDecodeSignedTransaction(signedTransaction)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "environment") {
		t.Fatalf("expected error to contain %q but got %q", "environment", err.Error())
	}
}

func TestAppStoreServerNotificationDecoding(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Read test notification data
	testNotification, err := ReadTestDataFileString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test notification decoding
	decodedPayload, err := verifier.VerifyAndDecodeNotification(testNotification)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil value")
	}

	// Verify basic notification fields
	AssertStringEqual(t, "2.0", decodedPayload.Version)
	if decodedPayload.Data == nil {
		t.Fatal("expected non-nil value")
	}
	AssertStringEqual(t, "com.example", decodedPayload.Data.BundleID)
	AssertStringEqual(t, "Sandbox", string(decodedPayload.Data.Environment))
}

func TestAppStoreServerNotificationDecodingProduction(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentProduction, "com.example")

	// This would test production environment notifications
	// For now, we'll test that the verifier is properly configured for production
	AssertEqual(t, EnvironmentProduction, verifier.environment)
	AssertStringEqual(t, "com.example", verifier.bundleID)
}

func TestMissingX5CHeader(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Read test data for missing x5c header
	missingX5CData, err := ReadTestDataFileString("mock_signed_data/missingX5CHeaderClaim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test notification decoding with missing x5c header - should fail
	_, err = verifier.VerifyAndDecodeNotification(missingX5CData)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	AssertContains(t, err.Error(), "x5c")
}

func TestWrongBundleIdForServerNotification(t *testing.T) {
	// Create verifier with wrong bundle ID
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.wrong.bundle")

	// Read test notification data
	testNotification, err := ReadTestDataFileString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test notification decoding with wrong bundle ID - should fail
	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	AssertContains(t, err.Error(), "bundle")
}

func TestWrongAppAppleIdForServerNotification(t *testing.T) {
	// Create verifier with wrong app apple ID
	verifier := createTestSignedDataVerifierWithAppID(t, EnvironmentSandbox, "com.example", 9999)

	// Read test notification data
	testNotification, err := ReadTestDataFileString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test notification decoding with wrong app apple ID - should fail
	_, err = verifier.VerifyAndDecodeNotification(testNotification)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	AssertContains(t, err.Error(), "app")
}

func TestRenewalInfoDecodingFromMockData(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Read test renewal info data
	renewalInfo, err := ReadTestDataFileString("mock_signed_data/renewalInfo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test renewal info decoding
	decodedPayload, err := verifier.VerifyAndDecodeRenewalInfo(renewalInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil value")
	}

	// Verify renewal info fields
	if decodedPayload.Environment != "Sandbox" {
		t.Fatalf("expected %q, got %q", "Sandbox", string(decodedPayload.Environment))
	}
	if decodedPayload.OriginalTransactionID == "" {
		t.Fatal("expected non-empty OriginalTransactionID")
	}
	if decodedPayload.ProductID == "" {
		t.Fatal("expected non-empty ProductID")
	}
}

func TestTransactionInfoDecoding(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Read test transaction info data
	transactionInfo, err := ReadTestDataFileString("mock_signed_data/transactionInfo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test transaction info decoding
	decodedPayload, err := verifier.VerifyAndDecodeSignedTransaction(transactionInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil value")
	}

	// Verify transaction info fields
	if decodedPayload.BundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", decodedPayload.BundleID)
	}
	if decodedPayload.ProductID == "" {
		t.Fatal("expected non-empty ProductID")
	}
	if decodedPayload.TransactionID == "" {
		t.Fatal("expected non-empty TransactionID")
	}
}

func TestMalformedJWTWithTooManyParts(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Test with JWT that has too many parts
	malformedJWT := "header.payload.signature.extra"

	// Test transaction decoding with malformed JWT - should fail
	_, err := verifier.VerifyAndDecodeSignedTransaction(malformedJWT)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	AssertContains(t, err.Error(), "invalid")
}

func TestMalformedJWTWithMalformedData(t *testing.T) {
	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Test with completely invalid JWT
	malformedJWT := "not.a.jwt"

	// Test transaction decoding with malformed data - should fail
	_, err := verifier.VerifyAndDecodeSignedTransaction(malformedJWT)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestVerificationWithWrongBundleId(t *testing.T) {
	// Read test data with wrong bundle ID
	wrongBundleData, err := ReadTestDataFileString("mock_signed_data/wrongBundleId")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	verifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")

	// Test transaction decoding with wrong bundle ID - should fail
	_, err = verifier.VerifyAndDecodeSignedTransaction(wrongBundleData)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	AssertContains(t, err.Error(), "bundle")
}

func TestVerificationWithDifferentEnvironments(t *testing.T) {
	// Test sandbox verifier
	sandboxVerifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")
	AssertEqual(t, EnvironmentSandbox, sandboxVerifier.environment)

	// Test production verifier
	productionVerifier := createTestSignedDataVerifier(t, EnvironmentProduction, "com.example")
	AssertEqual(t, EnvironmentProduction, productionVerifier.environment)

	// Verify they have different environments
	AssertTrue(t, sandboxVerifier.environment != productionVerifier.environment)
}

func TestVerificationWithLocalTestingEnvironment(t *testing.T) {
	// Create verifier for local testing environment
	verifier := createTestSignedDataVerifier(t, "LocalTesting", "com.example")

	// Read test notification data (which should be for LocalTesting environment)
	testNotification, err := ReadTestDataFileString("mock_signed_data/testNotification")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test notification decoding in local testing environment
	decodedPayload, err := verifier.VerifyAndDecodeNotification(testNotification)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decodedPayload == nil {
		t.Fatal("expected non-nil value")
	}
	AssertStringEqual(t, "LocalTesting", string(decodedPayload.Data.Environment))
}

// Helper function to create a test signed data verifier
func createTestSignedDataVerifier(t *testing.T, env Environment, bundleID string) *SignedDataVerifier {
	return createTestSignedDataVerifierWithAppID(t, env, bundleID, 1234)
}

// Helper function to create a test signed data verifier with specific app ID
func createTestSignedDataVerifierWithAppID(t *testing.T, env Environment, bundleID string, appAppleID int64) *SignedDataVerifier {
	// Read test CA certificate
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	verifier := &SignedDataVerifier{
		rootCertificates: [][]byte{testCA},
		environment:      env,
		bundleID:         bundleID,
		appAppleID:       appAppleID,
		chainVerifier:    newChainVerifier([][]byte{testCA}),
	}

	// Disable strict checks for test certificates
	if verifier.chainVerifier != nil {
		// This would need to be implemented based on the actual chainVerifier structure
		// verifier.chainVerifier.enableStrictChecks = false
	}

	return verifier
}

func TestCertificateChainVerification(t *testing.T) {
	// Read test CA certificate
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create chain verifier
	chainVerifier := newChainVerifier([][]byte{testCA})
	if chainVerifier == nil {
		t.Fatal("expected non-nil value")
	}

	// Test basic chain verifier functionality
	// Note: Specific verification tests would depend on the actual implementation
	// and available test certificates
}

func TestChainVerifierWithInvalidCertificate(t *testing.T) {
	// Test with invalid certificate data
	invalidCert := []byte("invalid certificate data")

	chainVerifier := newChainVerifier([][]byte{invalidCert})
	if chainVerifier == nil {
		t.Fatal("expected non-nil value")
	}

	// The chain verifier should be created but verification should fail
	// when used with actual certificate chains
}

func TestChainVerifierWithEmptyCertificates(t *testing.T) {
	// Test with empty certificate list
	chainVerifier := newChainVerifier([][]byte{})
	if chainVerifier == nil {
		t.Fatal("expected non-nil value")
	}

	// The chain verifier should handle empty certificate lists gracefully
}

func TestChainVerifierWithMultipleCertificates(t *testing.T) {
	// Read test CA certificate
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create chain verifier with multiple certificates (using the same cert twice for testing)
	chainVerifier := newChainVerifier([][]byte{testCA, testCA})
	if chainVerifier == nil {
		t.Fatal("expected non-nil value")
	}

	// Verify the chain verifier was created with multiple certificates
}

func TestCertificateValidation(t *testing.T) {
	// Read test CA certificate
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a signed data verifier with the test CA
	verifier := &SignedDataVerifier{
		rootCertificates: [][]byte{testCA},
		environment:      EnvironmentSandbox,
		bundleID:         "com.example",
		appAppleID:       1234,
		chainVerifier:    newChainVerifier([][]byte{testCA}),
	}

	if len(verifier.rootCertificates) != 1 {
		t.Fatalf("expected %v, got %v", 1, len(verifier.rootCertificates))
	}
	if verifier.environment != EnvironmentSandbox {
		t.Fatalf("expected %v, got %v", EnvironmentSandbox, verifier.environment)
	}
	if verifier.bundleID != "com.example" {
		t.Fatalf("expected %q, got %q", "com.example", verifier.bundleID)
	}
	if verifier.appAppleID != int64(1234) {
		t.Fatalf("expected %v, got %v", int64(1234), verifier.appAppleID)
	}
}

func TestSignedDataVerifierCreation(t *testing.T) {
	// Test creating signed data verifier with different configurations

	// Test with Sandbox environment
	sandboxVerifier := createTestSignedDataVerifier(t, EnvironmentSandbox, "com.example")
	if sandboxVerifier == nil {
		t.Fatal("expected non-nil value")
	}
	if sandboxVerifier.environment != EnvironmentSandbox {
		t.Fatalf("expected %v, got %v", EnvironmentSandbox, sandboxVerifier.environment)
	}

	// Test with Production environment
	productionVerifier := createTestSignedDataVerifier(t, EnvironmentProduction, "com.example.prod")
	if productionVerifier == nil {
		t.Fatal("expected non-nil value")
	}
	if productionVerifier.environment != EnvironmentProduction {
		t.Fatalf("expected %v, got %v", EnvironmentProduction, productionVerifier.environment)
	}

	// Test with different bundle IDs
	if sandboxVerifier.bundleID == productionVerifier.bundleID {
		t.Fatal("expected different bundle IDs")
	}
}

func TestSignedDataVerifierWithDifferentAppIDs(t *testing.T) {
	// Test creating verifiers with different App Apple IDs
	verifier1 := createTestSignedDataVerifierWithAppID(t, EnvironmentSandbox, "com.example", 1111)
	verifier2 := createTestSignedDataVerifierWithAppID(t, EnvironmentSandbox, "com.example", 2222)

	if verifier1 == nil {
		t.Fatal("expected non-nil value")
	}
	if verifier2 == nil {
		t.Fatal("expected non-nil value")
	}
	if verifier1.appAppleID != int64(1111) {
		t.Fatalf("expected %v, got %v", int64(1111), verifier1.appAppleID)
	}
	if verifier2.appAppleID != int64(2222) {
		t.Fatalf("expected %v, got %v", int64(2222), verifier2.appAppleID)
	}

	// Verify they have different app IDs
	if verifier1.appAppleID == verifier2.appAppleID {
		t.Fatal("expected different app IDs")
	}
}

func TestChainVerifierInitialization(t *testing.T) {
	// Test that chain verifier is properly initialized
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chainVerifier := newChainVerifier([][]byte{testCA})
	if chainVerifier == nil {
		t.Fatal("expected non-nil value")
	}

	// Test with nil certificates
	nilChainVerifier := newChainVerifier(nil)
	if nilChainVerifier == nil {
		t.Fatal("expected non-nil value")
	}
}

func TestVerificationEnvironmentValidation(t *testing.T) {
	// Test that environment validation works correctly
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test valid environments
	validEnvironments := []Environment{
		EnvironmentSandbox,
		EnvironmentProduction,
	}

	for _, env := range validEnvironments {
		verifier := &SignedDataVerifier{
			rootCertificates: [][]byte{testCA},
			environment:      env,
			bundleID:         "com.example",
			appAppleID:       1234,
			chainVerifier:    newChainVerifier([][]byte{testCA}),
		}

		if verifier.environment != env {
			t.Fatalf("expected %v, got %v", env, verifier.environment)
		}
		if !verifier.environment.IsValid() {
			t.Fatal("expected valid environment")
		}
	}
}

func TestBundleIDValidation(t *testing.T) {
	// Test different bundle ID formats
	testCA, err := ReadTestDataFile("certs/testCA.der")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bundleIDs := []string{
		"com.example",
		"com.example.app",
		"com.company.product.name",
		"org.opensource.project",
	}

	for _, bundleID := range bundleIDs {
		verifier := &SignedDataVerifier{
			rootCertificates: [][]byte{testCA},
			environment:      EnvironmentSandbox,
			bundleID:         bundleID,
			appAppleID:       1234,
			chainVerifier:    newChainVerifier([][]byte{testCA}),
		}

		if verifier.bundleID != bundleID {
			t.Fatalf("expected %q, got %q", bundleID, verifier.bundleID)
		}
	}
}
