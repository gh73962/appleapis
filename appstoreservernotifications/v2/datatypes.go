package appstoreservernotifications

import "time"

// ResponseBody contains the version 2 notification data.
// See https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2
type ResponseBody struct {
	SignedPayload string `json:"signedPayload"`
}

// DecodedPayload contains the version 2 notification data.
// See https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
type DecodedPayload struct {
	NotificationType      NotificationType       `json:"notificationType,omitempty"`
	Subtype               Subtype                `json:"subtype,omitempty"`
	NotificationUUID      string                 `json:"notificationUUID,omitempty"`
	Data                  *Data                  `json:"data,omitempty"`
	Version               string                 `json:"version,omitempty"`
	SignedDate            int64                  `json:"signedDate,omitempty"`
	Summary               *Summary               `json:"summary,omitempty"`
	ExternalPurchaseToken *ExternalPurchaseToken `json:"externalPurchaseToken,omitempty"`
}

func (d *DecodedPayload) GetSignedDate() time.Time {
	return time.UnixMilli(d.SignedDate)
}

// Data contains the app metadata and signed renewal and transaction information.
// See https://developer.apple.com/documentation/appstoreservernotifications/data
type Data struct {
	Environment           string `json:"environment,omitempty"`
	AppAppleID            int64  `json:"appAppleId,omitempty"`
	BundleID              string `json:"bundleId,omitempty"`
	BundleVersion         string `json:"bundleVersion,omitempty"`
	SignedRenewalInfo     string `json:"signedRenewalInfo,omitempty"`
	SignedTransactionInfo string `json:"signedTransactionInfo,omitempty"`
	Status                int    `json:"status,omitempty"`
}

// IsActive returns true if the auto-renewable subscription is active.
func (d *Data) IsActive() bool {
	return d.Status == 1
}

// IsExpired returns true if the auto-renewable subscription is expired.
func (d *Data) IsExpired() bool {
	return d.Status == 2
}

// IsInBillingRetry returns true if the auto-renewable subscription is in a billing retry period.
func (d *Data) IsInBillingRetry() bool {
	return d.Status == 3
}

// IsInBillingGracePeriod returns true if the auto-renewable subscription is in a Billing Grace Period.
func (d *Data) IsInBillingGracePeriod() bool {
	return d.Status == 4
}

// IsRevoked returns true if the auto-renewable subscription is revoked.
func (d *Data) IsRevoked() bool {
	return d.Status == 5
}

func (d *Data) IsSandbox() bool {
	return d.Environment == "Sandbox"
}

func (d *Data) IsProduction() bool {
	return d.Environment == "Production"
}

// Summary contains the summary data for subscription renewal date extension notifications.
// See https://developer.apple.com/documentation/appstoreservernotifications/summary
type Summary struct {
	Environment            string   `json:"environment,omitempty"`
	AppAppleID             int64    `json:"appAppleId,omitempty"`
	BundleID               string   `json:"bundleId,omitempty"`
	ProductID              string   `json:"productId,omitempty"`
	RequestIdentifier      string   `json:"requestIdentifier,omitempty"`
	StorefrontCountryCodes []string `json:"storefrontCountryCodes,omitempty"`
	SucceededCount         int64    `json:"succeededCount,omitempty"`
	FailedCount            int64    `json:"failedCount,omitempty"`
}

func (s *Summary) IsSandbox() bool {
	return s.Environment == "Sandbox"
}

func (s *Summary) IsProduction() bool {
	return s.Environment == "Production"
}

// ExternalPurchaseToken contains external purchase token information.
// See https://developer.apple.com/documentation/appstoreservernotifications/externalpurchasetoken
type ExternalPurchaseToken struct {
	ExternalPurchaseID  string `json:"externalPurchaseId,omitempty"`
	TokenCreationDate   int64  `json:"tokenCreationDate,omitempty"`
	AppAppleID          int64  `json:"appAppleId,omitempty"`
	BundleID            string `json:"bundleId,omitempty"`
	TokenType           string `json:"tokenType,omitempty"`
	TokenExpirationDate int64  `json:"tokenExpirationDate,omitempty"`
}

func (e *ExternalPurchaseToken) GetTokenCreationDate() time.Time {
	return time.UnixMilli(e.TokenCreationDate)
}

func (e *ExternalPurchaseToken) GetTokenExpirationDate() time.Time {
	return time.UnixMilli(e.TokenExpirationDate)
}

func (e *ExternalPurchaseToken) IsInitialAcquisition() bool {
	return e.TokenType == "ACQUISITION"
}

func (e *ExternalPurchaseToken) IsServices() bool {
	return e.TokenType == "SERVICES"
}
