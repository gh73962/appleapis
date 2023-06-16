// Package datatypes See https://developer.apple.com/documentation/appstoreserverapi/data_types
package datatypes

import (
	"encoding/json"
	"time"
)

const (
	BasePath        = "https://api.storekit.itunes.apple.com/inApps/v1/"
	SandboxBasePath = "https://api.storekit-sandbox.itunes.apple.com/inApps/v1/"
)

// OfferType see https://developer.apple.com/documentation/appstoreserverapi/offertype
type OfferType int

const (
	IntroductoryOffer              OfferType = 1
	PromotionalOffer               OfferType = 2
	OfferWithSubscriptionOfferCode OfferType = 3
)

// Environment see https://developer.apple.com/documentation/appstoreserverapi/environment
type Environment string

const (
	Sandbox    Environment = "Sandbox"
	Production Environment = "Production"
)

// InAppOwnershipType see https://developer.apple.com/documentation/appstoreserverapi/inappownershiptype
type InAppOwnershipType string

const (
	FamilyShared InAppOwnershipType = "FAMILY_SHARED"
	Purchased    InAppOwnershipType = "PURCHASED"
)

// TransactionType see https://developer.apple.com/documentation/appstoreserverapi/type
type TransactionType string

const (
	AutoRenewableSubscription TransactionType = "Auto-Renewable Subscription"
	NonConsumable             TransactionType = "Non-Consumable"
	Consumable                TransactionType = "Consumable"
	NonRenewingSubscription   TransactionType = "Non-Renewing Subscription"
)

// JWSDecodedHeader https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type JWSDecodedHeader struct {
	Alg string   `json:"alg,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

// JWSTransaction see https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
type JWSTransaction struct {
	Header    JWSDecodedHeader
	Payload   JWSTransactionDecodedPayload
	Signature string
}

// JWSTransactionDecodedPayload https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type JWSTransactionDecodedPayload struct {
	AppAccountToken             string             `json:"appAccountToken,omitempty"`
	BundleID                    string             `json:"bundleId,omitempty"`
	Environment                 Environment        `json:"environment,omitempty"`
	ExpiresDate                 int64              `json:"expiresDate,omitempty"`
	InAppOwnershipType          InAppOwnershipType `json:"inAppOwnershipType,omitempty"`
	IsUpgraded                  bool               `json:"isUpgraded,omitempty"`
	OfferIdentifier             string             `json:"offerIdentifier,omitempty"`
	OfferType                   OfferType          `json:"offerType,omitempty"`
	OriginalPurchaseDate        int64              `json:"originalPurchaseDate,omitempty"`
	OriginalTransactionID       string             `json:"originalTransactionId,omitempty"`
	ProductID                   string             `json:"productId,omitempty"`
	PurchaseDate                int64              `json:"purchaseDate,omitempty"`
	Quantity                    int                `json:"quantity,omitempty"`
	RevocationDate              int64              `json:"revocationDate,omitempty"`
	RevocationReason            int                `json:"revocationReason,omitempty"`
	SignedDate                  int64              `json:"signedDate,omitempty"`
	Storefront                  string             `json:"storefront,omitempty"`
	StorefrontID                string             `json:"storefrontId,omitempty"`
	SubscriptionGroupIdentifier string             `json:"subscriptionGroupIdentifier,omitempty"`
	TransactionID               string             `json:"transactionId,omitempty"`
	TransactionReason           string             `json:"transactionReason,omitempty"`
	Type                        TransactionType    `json:"type,omitempty"`
	WebOrderLineItemID          string             `json:"webOrderLineItemId,omitempty"`
}

func (j *JWSTransactionDecodedPayload) GetPurchaseTime() time.Time {
	return time.Unix(j.PurchaseDate/1e3, 0)
}

func (j *JWSTransactionDecodedPayload) GetExpiresTime() time.Time {
	return time.Unix(j.ExpiresDate/1e3, 0)
}

func (j *JWSTransactionDecodedPayload) GetSignedTime() time.Time {
	return time.Unix(j.SignedDate/1e3, 0)
}

// JWSRenewalInfo see https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfo
type JWSRenewalInfo struct {
	Header    JWSDecodedHeader
	Payload   JWSRenewalInfoDecodedPayload
	Signature string
}

// ExpirationIntent see https://developer.apple.com/documentation/appstoreserverapi/expirationintent
type ExpirationIntent int

const (
	CanceledSubscription    ExpirationIntent = 1
	BillingError            ExpirationIntent = 2
	NotConsentPriceIncrease ExpirationIntent = 3
	NotAvailable            ExpirationIntent = 4
	OtherReason             ExpirationIntent = 5
)

// JWSRenewalInfoDecodedPayload https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	AutoRenewProductID          string           `json:"autoRenewProductId,omitempty"`
	AutoRenewStatus             int              `json:"autoRenewStatus,omitempty"`
	Environment                 Environment      `json:"environment,omitempty"`
	ExpirationIntent            ExpirationIntent `json:"expirationIntent,omitempty"`
	GracePeriodExpiresDate      int64            `json:"gracePeriodExpiresDate,omitempty"`
	IsInBillingRetryPeriod      bool             `json:"isInBillingRetryPeriod,omitempty"`
	OfferIdentifier             string           `json:"offerIdentifier,omitempty"`
	OfferType                   OfferType        `json:"offerType,omitempty"`
	OriginalTransactionID       string           `json:"originalTransactionId,omitempty"`
	PriceIncreaseStatus         int              `json:"priceIncreaseStatus,omitempty"`
	ProductID                   string           `json:"productId,omitempty"`
	RecentSubscriptionStartDate int64            `json:"recentSubscriptionStartDate,omitempty"`
	RenewalDate                 int64            `json:"renewalDate,omitempty"`
	SignedDate                  int64            `json:"signedDate,omitempty"`
}

func (j *JWSRenewalInfoDecodedPayload) GetRenewalTime() time.Time {
	return time.Unix(j.RenewalDate/1e3, 0)
}

func (j *JWSRenewalInfoDecodedPayload) IsAutoRenew() bool {
	return j.AutoRenewStatus == 1
}

func (j *JWSRenewalInfoDecodedPayload) IsConsentedPriceIncrease() bool {
	return j.PriceIncreaseStatus == 1
}

// TransactionInfoResponse see https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type TransactionInfoResponse struct {
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// StatusResponse https://developer.apple.com/documentation/appstoreserverapi/statusresponse
type StatusResponse struct {
	BundleID    string                            `json:"bundleId,omitempty"`
	AppAppleID  int64                             `json:"appAppleId,omitempty"`
	Environment string                            `json:"environment,omitempty"`
	Data        []SubscriptionGroupIdentifierItem `json:"data,omitempty"`
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier,omitempty"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions,omitempty"`
}

// SubscriptionStatus see https://developer.apple.com/documentation/appstoreserverapi/status
type SubscriptionStatus int

const (
	Active             SubscriptionStatus = 1
	Expired            SubscriptionStatus = 2
	BillingRetryPeriod SubscriptionStatus = 3
	BillingGracePeriod SubscriptionStatus = 4
	Revoked            SubscriptionStatus = 5
)

type LastTransactionsItem struct {
	OriginalTransactionID string             `json:"originalTransactionId,omitempty"`
	Status                SubscriptionStatus `json:"status,omitempty"`
	SignedRenewalInfo     string             `json:"signedRenewalInfo,omitempty"`
	SignedTransactionInfo string             `json:"signedTransactionInfo,omitempty"`
}

// HistoryResponse see https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	Revision           string      `json:"revision,omitempty"`
	BundleID           string      `json:"bundleId,omitempty"`
	AppAppleID         int         `json:"appAppleId,omitempty"`
	Environment        Environment `json:"environment,omitempty"`
	HasMore            bool        `json:"hasMore,omitempty"`
	SignedTransactions []string    `json:"signedTransactions,omitempty"`
}

// OrderLookupResponse see https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             int      `json:"status,omitempty"`
	SignedTransactions []string `json:"signedTransactions,omitempty"`
}

func (o *OrderLookupResponse) IsValid() bool {
	return o.Status == 0
}

// RefundHistoryResponse see https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundHistoryResponse struct {
	Revision           string   `json:"revision,omitempty"`
	HasMore            bool     `json:"hasMore,omitempty"`
	SignedTransactions []string `json:"signedTransactions,omitempty"`
}

// ErrorResponse See https://developer.apple.com/documentation/appstoreserverapi/error_codes
type ErrorResponse struct {
	HTTPStatus   int    `json:"httpStatus,omitempty"`
	ErrorCode    int64  `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e == nil {
		return ""
	}
	data, _ := json.Marshal(e)
	return string(data)
}

// ConsumptionRequest see https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequest struct {
	AccountTenure            int    `json:"accountTenure,omitempty"`
	AppAccountToken          string `json:"appAccountToken,omitempty"`
	ConsumptionStatus        int    `json:"consumptionStatus,omitempty"`
	CustomerConsented        bool   `json:"customerConsented,omitempty"`
	DeliveryStatus           int    `json:"deliveryStatus,omitempty"`
	LifetimeDollarsPurchased int    `json:"lifetimeDollarsPurchased,omitempty"`
	Platform                 int    `json:"platform,omitempty"`
	PlayTime                 int    `json:"playTime,omitempty"`
	SampleContentProvided    bool   `json:"sampleContentProvided,omitempty"`
	UserStatus               int    `json:"userStatus,omitempty"`
}

// NotificationHistoryRequest see https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	StartDate           int64  `json:"startDate,omitempty"`
	EndDate             int64  `json:"endDate,omitempty"`
	NotificationType    string `json:"notificationType,omitempty"`
	NotificationSubtype string `json:"notificationSubtype,omitempty"`
	OnlyFailures        bool   `json:"onlyFailures,omitempty"`
	TransactionID       string `json:"transactionId,omitempty"`
}

// NotificationHistoryResponse see https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponse struct {
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory,omitempty"`
	HasMore             bool                              `json:"hasMore,omitempty"`
	PaginationToken     string                            `json:"paginationToken,omitempty"`
}

type NotificationHistoryResponseItem struct {
	SendAttempts  []SendAttemptItem `json:"sendAttempts,omitempty"`
	SignedPayload string            `json:"signedPayload,omitempty"`
}

type SendAttemptItem struct {
	AttemptDate       int64  `json:"attemptDate,omitempty"`
	SendAttemptResult string `json:"sendAttemptResult,omitempty"`
}

type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken,omitempty"`
}
