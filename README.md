# Apple APIs Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/gh73962/appleapis.svg)](https://pkg.go.dev/github.com/gh73962/appleapis)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go client library for Apple APIs, including the App Store Server API and Server Notifications.

## Features

- 🚀 **App Store Server API v1**: Complete implementation of Apple's App Store Server API
- 📱 **Server Notifications v2**: Support for App Store Server Notifications 
- 🔐 **JWT Authentication**: Automatic JWT token generation and signing
- ✅ **JWS Verification**: Built-in verification of signed data from Apple
- 🌍 **Multi-Environment**: Support for Production, Sandbox, and Local Testing
- 📊 **Type Safety**: Full Go struct definitions for all API responses
- 🔄 **Context Support**: Context-aware API calls for better control
- 📦 **Zero Dependencies**: Minimal external dependencies

## Installation

```bash
go get github.com/gh73962/appleapis
```

## Quick Start

### 1. Configuration

First, create a client configuration with your App Store Connect credentials:

```go
package main

import (
    "github.com/gh73962/appleapis/appstoreserver/v1"
)

func main() {
   opts := []Option{
		WithAppAppleID(1234),
		WithBundleID("com.example"),
		WithEnvironment(EnvironmentProduction),
		WithKeyID("your_key_id"),
		WithIssuerID("your_issuer_id"),
		WithPrivateKey(your_pk),
		WithRootCertificates([][]byte{your_cert}),
        WithEnableAutoDecode(), // Recommended to be true for most use cases.
	}

	client, err := New(opts...)
	if err != nil {
		return nil, err
	}
}
```

### 2. Get Transaction History

```go
ctx := context.Background()

req := &TransactionHistoryRequest{
    TransactionID:                "1234",
    ProductIDs:                   []string{"com.example.1", "com.example.2"},
    ProductTypes:                 []ProductType{TypeConsumable, TypeAutoRenewableSubscription},
    StartDate:                    time.Now(),
    EndDate:                      time.Now(),
    SubscriptionGroupIdentifiers: []string{"sub_group_id", "sub_group_id_2"},
    InAppOwnershipType:           InAppOwnershipTypeFamilyShared,
    Revoked:                      false,
    Revision:                     "revision_input",
}
req.SetSortASC()

response, err := client.GetTransactionHistory(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d transactions\n", len(response.Payloads))
```

### 3. Get Transaction Info

```go
transactionInfo, err := client.GetTransactionInfo(ctx, "1000000123456789")
if err != nil {
    // handle error
}

fmt.Printf("Payload: %+v\n", transactionInfo.Payload)
```

### 4. Get Subscription Status

```go
status, err := client.GetAllSubscriptionStatuses(ctx, "1000000123456789")
if err != nil {
    // handle error
}

for _, data := range status.Data {
    fmt.Printf("Bundle ID: %s\n", data.BundleID)
    for _, transaction := range data.LastTransactions {
        fmt.Printf("Status: %s\n", transaction.Status)
    }
}
```

### 4. Easily SendConsumptionInfo

```go
var req appstoreserver.ConsumptionRequest
req.SetSetLifetimeDollarsPurchased(5.99)
req.SetLifetimeDollarsRefunded(5.99)
req.SetAccountTenure(365)
req.SetPlayTime(96*time.Hour)

if err := client.SendConsumptionInfo(ctx, &req); err != nil {
    // handle error
}
```

## Server Notifications

Handle App Store Server Notifications v2:
Ref [receiving-app-store-server-notifications](https://developer.apple.com/documentation/appstoreservernotifications/receiving-app-store-server-notifications)
```go
import "github.com/gh73962/appleapis/appstoreservernotifications/v2"

// In your webhook handler
func handleNotification(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusBadRequest)
        return
    }

    var notification appstoreservernotifications.ResponseBody
    if err := json.Unmarshal(body, &notification); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    decodedPayload, err := client.Verifier.VerifyAndDecodeNotification(notification.SignedPayload)
    if err != nil {
        http.Error(w, "VerifyAndDecodeNotification", http.StatusBadRequest)
        return
    }

    switch decodedPayload.NotificationType{
        case appstoreservernotifications.TypeConsumptionRequest:
        // Do 
        case appstoreservernotifications.TypeRefund
        // Do
    }
    w.WriteHeader(http.StatusOK)
}
```

## API Coverage

### App Store Server API v1
Web Service Endpoint list on [Apple App Store Server API Documentation](https://developer.apple.com/documentation/appstoreserverapi)
- ✅ Get Transaction History V2
- ✅ Get Transaction Info  
- ✅ Get All Subscription Statuses
- ✅ Set App Account Token
- ✅ Send Consumption Information
- ✅ Look Up Order ID
- ✅ Get Refund History V2
- ✅ Extend Subscription Renewal Date
- ✅ Extend Subscription Renewal Dates for All Active Subscribers
- ✅ Get Notification History
- ✅ Request Test Notification
- ✅ Get Test Notification Status


### Server Notifications v2

- ✅ All notification types supported
- ✅ Structured data types for all payloads
- ✅ Built-in JWS signature verification

## Configuration Options

| Option | Description | Required |
|--------|-------------|----------|
| `PrivateKey` | Your private key from App Store Connect (.p8 file content) | Yes |
| `KeyID` | Key ID from App Store Connect | Yes |
| `IssuerID` | Issuer ID from App Store Connect | Yes |
| `BundleID` | Your app's bundle identifier | Yes |
| `Environment` | `EnvironmentSandbox` or `EnvironmentProduction` | Yes |
| `AppAppleID` | Your app's Apple ID | On Production |
| `RootCertificates` | Custom Apple Root CA certificates | No |
| `EnableOnlineChecks` | Enable online certificate verification | No |
| `HTTPClient` | Custom HTTP client | No |


## Testing

Run the test suite:

```bash
go test ./...
```


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feature: your desc'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Apple App Store Server API Documentation](https://developer.apple.com/documentation/appstoreserverapi)
- [Apple App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications)

## Related Projects

- [app-store-server-library-python](https://github.com/apple/app-store-server-library-python) - Official Python library from Apple
