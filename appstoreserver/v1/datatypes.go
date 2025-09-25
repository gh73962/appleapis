package appstoreserver

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/gh73962/appleapis/appstoreservernotifications/v2"
)

// JWSTransactionDecodedPayload contains transaction information signed by the App Store.
// See https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type JWSTransactionDecodedPayload struct {
	AppAccountToken             *string                `json:"appAccountToken,omitempty"`
	BundleID                    string                 `json:"bundleId"`
	Environment                 Environment            `json:"environment"`
	ExpiresDate                 *int64                 `json:"expiresDate,omitempty"`
	InAppOwnershipType          InAppOwnershipType     `json:"inAppOwnershipType"`
	IsUpgraded                  *bool                  `json:"isUpgraded,omitempty"`
	OfferIdentifier             *string                `json:"offerIdentifier,omitempty"`
	OfferType                   *SubscriptionOfferType `json:"offerType,omitempty"`
	OriginalPurchaseDate        int64                  `json:"originalPurchaseDate"`
	OriginalTransactionID       string                 `json:"originalTransactionId"`
	ProductID                   string                 `json:"productId"`
	PurchaseDate                int64                  `json:"purchaseDate"`
	Quantity                    int                    `json:"quantity"`
	RevocationDate              *int64                 `json:"revocationDate,omitempty"`
	RevocationReason            *RevocationReason      `json:"revocationReason,omitempty"`
	SignedDate                  int64                  `json:"signedDate"`
	SubscriptionGroupIdentifier *string                `json:"subscriptionGroupIdentifier,omitempty"`
	TransactionID               string                 `json:"transactionId"`
	TransactionReason           *TransactionReason     `json:"transactionReason,omitempty"`
	Type                        ProductType            `json:"type"`
	WebOrderLineItemID          *string                `json:"webOrderLineItemId,omitempty"`
}

// GetExpiresDate returns the expiration date as a time.Time
func (t *JWSTransactionDecodedPayload) GetExpiresDate() *time.Time {
	if t.ExpiresDate == nil {
		return nil
	}
	timestamp := time.Unix(*t.ExpiresDate/1000, 0)
	return &timestamp
}

// GetOriginalPurchaseDate returns the original purchase date as a time.Time
func (t *JWSTransactionDecodedPayload) GetOriginalPurchaseDate() time.Time {
	return time.Unix(t.OriginalPurchaseDate/1000, 0)
}

// GetPurchaseDate returns the purchase date as a time.Time
func (t *JWSTransactionDecodedPayload) GetPurchaseDate() time.Time {
	return time.Unix(t.PurchaseDate/1000, 0)
}

// GetRevocationDate returns the revocation date as a time.Time
func (t *JWSTransactionDecodedPayload) GetRevocationDate() *time.Time {
	if t.RevocationDate == nil {
		return nil
	}
	timestamp := time.Unix(*t.RevocationDate/1000, 0)
	return &timestamp
}

// GetSignedDate returns the signed date as a time.Time
func (t *JWSTransactionDecodedPayload) GetSignedDate() time.Time {
	return time.Unix(t.SignedDate/1000, 0)
}

// JWSRenewalInfoDecodedPayload contains subscription renewal information signed by the App Store.
// See https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	AutoRenewProductID          *string                `json:"autoRenewProductId,omitempty"`
	AutoRenewStatus             AutoRenewStatus        `json:"autoRenewStatus"`
	Environment                 Environment            `json:"environment"`
	ExpirationIntent            *ExpirationIntent      `json:"expirationIntent,omitempty"`
	GracePeriodExpiresDate      *int64                 `json:"gracePeriodExpiresDate,omitempty"`
	IsInBillingRetryPeriod      *bool                  `json:"isInBillingRetryPeriod,omitempty"`
	OfferIdentifier             *string                `json:"offerIdentifier,omitempty"`
	OfferType                   *SubscriptionOfferType `json:"offerType,omitempty"`
	OriginalTransactionID       string                 `json:"originalTransactionId"`
	PriceIncreaseStatus         *int                   `json:"priceIncreaseStatus,omitempty"`
	ProductID                   string                 `json:"productId"`
	RecentSubscriptionStartDate *int64                 `json:"recentSubscriptionStartDate,omitempty"`
	SignedDate                  int64                  `json:"signedDate"`
}

// GetGracePeriodExpiresDate returns the grace period expiration date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetGracePeriodExpiresDate() *time.Time {
	if r.GracePeriodExpiresDate == nil {
		return nil
	}
	timestamp := time.Unix(*r.GracePeriodExpiresDate/1000, 0)
	return &timestamp
}

// GetRecentSubscriptionStartDate returns the recent subscription start date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetRecentSubscriptionStartDate() *time.Time {
	if r.RecentSubscriptionStartDate == nil {
		return nil
	}
	timestamp := time.Unix(*r.RecentSubscriptionStartDate/1000, 0)
	return &timestamp
}

// GetSignedDate returns the signed date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetSignedDate() time.Time {
	return time.Unix(r.SignedDate/1000, 0)
}

// ConsumptionRequest contains consumption information for a transaction.
// See https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequest struct {
	TransactionID            string                   `json:"-"`
	AccountTenure            AccountTenure            `json:"accountTenure"`
	AppAccountToken          string                   `json:"appAccountToken"`
	ConsumptionStatus        ConsumptionStatus        `json:"consumptionStatus"`
	CustomerConsented        bool                     `json:"customerConsented"`
	DeliveryStatus           DeliveryStatus           `json:"deliveryStatus"`
	LifetimeDollarsPurchased LifetimeDollarsPurchased `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded  LifetimeDollarsRefunded  `json:"lifetimeDollarsRefunded"`
	Platform                 Platform                 `json:"platform"`
	PlayTime                 PlayTime                 `json:"playTime"`
	RefundPreference         RefundPreference         `json:"refundPreference,omitempty"`
	SampleContentProvided    bool                     `json:"sampleContentProvided"`
	UserStatus               UserStatus               `json:"userStatus"`
}

func (c *ConsumptionRequest) Validate() error {
	if c.TransactionID == "" {
		return fmt.Errorf("transactionID is required")
	}
	if len(c.AppAccountToken) == 0 {
		return fmt.Errorf("appAccountToken is required")
	}

	return nil
}

// ExtendRenewalDateRequest contains information for extending a subscription renewal date.
// See https://developer.apple.com/documentation/appstoreserverapi/extendrenewalDaterequest
type ExtendRenewalDateRequest struct {
	OriginalTransactionID string           `json:"-"`
	ExtendByDays          int              `json:"extendByDays"`
	ExtendReasonCode      ExtendReasonCode `json:"extendReasonCode"`
	RequestIdentifier     string           `json:"requestIdentifier"`
}

func (e *ExtendRenewalDateRequest) Validate() error {
	if e.OriginalTransactionID == "" {
		return fmt.Errorf("originalTransactionID is required")
	}

	if len(e.RequestIdentifier) == 0 || len(e.RequestIdentifier) > 128 {
		return fmt.Errorf("requestIdentifier must be between 1 and 128 characters")
	}
	if e.ExtendByDays < 1 || e.ExtendByDays > 90 {
		return fmt.Errorf("extendByDays must be between 1 and 90")
	}
	if e.ExtendReasonCode < ExtendReasonCodeCustomerSatisfy || e.ExtendReasonCode > ExtendReasonCodeServiceIssue {
		return fmt.Errorf("extendReasonCode must be a valid value")
	}
	return nil
}

// MassExtendRenewalDateRequest contains information for extending subscription renewal dates for multiple users.
// See https://developer.apple.com/documentation/appstoreserverapi/massextendrenewalDaterequest
type MassExtendRenewalDateRequest struct {
	RequestIdentifier      string           `json:"requestIdentifier"`
	ExtendByDays           int              `json:"extendByDays"`
	ExtendReasonCode       ExtendReasonCode `json:"extendReasonCode"`
	ProductID              string           `json:"productId"`
	StorefrontCountryCodes []string         `json:"storefrontCountryCodes,omitempty"`
}

func (m *MassExtendRenewalDateRequest) Validate() error {
	if len(m.RequestIdentifier) == 0 || len(m.RequestIdentifier) > 128 {
		return fmt.Errorf("requestIdentifier must be between 1 and 128 characters")
	}
	if m.ExtendByDays < 1 || m.ExtendByDays > 90 {
		return fmt.Errorf("extendByDays must be between 1 and 90")
	}
	if m.ExtendReasonCode < ExtendReasonCodeCustomerSatisfy || m.ExtendReasonCode > ExtendReasonCodeServiceIssue {
		return fmt.Errorf("extendReasonCode must be a valid value")
	}
	if len(m.ProductID) == 0 {
		return fmt.Errorf("productID is required")
	}
	return nil
}

// NotificationHistoryRequest contains information for requesting notification history.
// See https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	StartTime           time.Time                                    `json:"-"`
	EndTime             time.Time                                    `json:"-"`
	StartDate           int64                                        `json:"startDate"`
	EndDate             int64                                        `json:"endDate"`
	NotificationType    appstoreservernotifications.NotificationType `json:"notificationType,omitempty"`
	NotificationSubtype appstoreservernotifications.Subtype          `json:"notificationSubtype,omitempty"`
	OnlyFailures        bool                                         `json:"onlyFailures,omitempty"`
	TransactionID       string                                       `json:"transactionId,omitempty"`
	PaginationToken     string                                       `json:"-"`
}

func (n *NotificationHistoryRequest) Validate() error {
	if !n.StartTime.IsZero() {
		n.StartDate = n.StartTime.UnixMilli()
	}
	if !n.EndTime.IsZero() {
		n.EndDate = n.EndTime.UnixMilli()
	}
	if n.StartDate <= 0 {
		return fmt.Errorf("startDate must be a valid timestamp in milliseconds")
	}
	if n.EndDate <= 0 {
		return fmt.Errorf("endDate must be a valid timestamp in milliseconds")
	}
	if n.EndDate < n.StartDate {
		return fmt.Errorf("endDate cannot be earlier than startDate")
	}
	return nil
}

// TransactionHistoryRequest contains information for requesting transaction history.
// See https://developer.apple.com/documentation/appstoreserverapi/get-transaction-history
type TransactionHistoryRequest struct {
	TransactionID                string
	Sort                         string
	ProductIDs                   []string
	ProductTypes                 []ProductType
	StartDate                    time.Time
	EndDate                      time.Time
	SubscriptionGroupIdentifiers []string
	InAppOwnershipType           InAppOwnershipType
	Revoked                      bool
	Revision                     string
}

func (t *TransactionHistoryRequest) SetSortASC() {
	t.Sort = "ASCENDING"
}

func (t *TransactionHistoryRequest) SetSortDESC() {
	t.Sort = "DESCENDING"
}

func (t *TransactionHistoryRequest) makeQuery() url.Values {
	queryParams := make(url.Values)
	if t.Sort != "" {
		queryParams.Add("sort", t.Sort)
	}
	for _, productID := range t.ProductIDs {
		queryParams.Add("productId", productID)
	}
	for _, productType := range t.ProductTypes {
		queryParams.Add("productType", string(productType))
	}
	if !t.StartDate.IsZero() {
		queryParams.Add("startDate", strconv.FormatInt(t.StartDate.UnixMilli(), 10))
	}
	if !t.EndDate.IsZero() {
		queryParams.Add("endDate", strconv.FormatInt(t.EndDate.UnixMilli(), 10))
	}
	for _, id := range t.SubscriptionGroupIdentifiers {
		queryParams.Add("subscriptionGroupIdentifier", id)
	}
	if t.InAppOwnershipType != "" {
		queryParams.Add("inAppOwnershipType", string(t.InAppOwnershipType))
	}
	if t.Revoked {
		queryParams.Add("revoked", strconv.FormatBool(t.Revoked))
	}
	if t.Revision != "" {
		queryParams.Add("revision", t.Revision)
	}
	return queryParams
}

// UpdateAppAccountTokenRequest contains information for updating an app account token.
// See https://developer.apple.com/documentation/appstoreserverapi/updateappaccounttokenrequest
type UpdateAppAccountTokenRequest struct {
	OriginalTransactionID string `json:"-"`
	AppAccountToken       string `json:"appAccountToken"`
}

func (u *UpdateAppAccountTokenRequest) Validate() error {
	if u.OriginalTransactionID == "" {
		return fmt.Errorf("originalTransactionID is required")
	}
	if len(u.AppAccountToken) == 0 {
		return fmt.Errorf("appAccountToken is required")
	}
	return nil
}

// TransactionInfoResponse A response that contains signed transaction information for a single transaction.
// See https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type TransactionInfoResponse struct {
	SignedTransactionInfo string                        `json:"signedTransactionInfo"`
	Payload               *JWSTransactionDecodedPayload `json:"-"`
}

// HistoryResponse contains the customer's transaction history for an app.
// See https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	Revision           string                          `json:"revision,omitempty"`
	HasMore            bool                            `json:"hasMore,omitempty"`
	BundleID           string                          `json:"bundleId,omitempty"`
	AppAppleID         int                             `json:"appAppleId,omitempty"`
	Environment        Environment                     `json:"environment,omitempty"`
	SignedTransactions []string                        `json:"signedTransactions,omitempty"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// OrderLookupResponse includes the order lookup status and an array of signed transactions for the in-app purchases in the order.
// See https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             OrderLookupStatus               `json:"status,omitempty"`
	SignedTransactions []string                        `json:"signedTransactions,omitempty"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// RefundHistoryResponse contains an array of signed JSON Web Signature (JWS) refunded transactions, and paging information.
// See https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundHistoryResponse struct {
	SignedTransactions []string                        `json:"signedTransactions,omitempty"`
	Revision           string                          `json:"revision,omitempty"`
	HasMore            bool                            `json:"hasMore,omitempty"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// ExtendRenewalDateResponse indicates whether an individual renewal-date extension succeeded, and related details.
// See https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldateresponse
type ExtendRenewalDateResponse struct {
	OriginalTransactionID string `json:"originalTransactionId,omitempty"`
	WebOrderLineItemID    string `json:"webOrderLineItemId,omitempty"`
	Success               bool   `json:"success,omitempty"`
	EffectiveDate         int64  `json:"effectiveDate,omitempty"`
}

// MassExtendRenewalDateResponse indicates the server successfully received the subscription-renewal-date extension request.
// See https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldateresponse
type MassExtendRenewalDateResponse struct {
	RequestIdentifier string `json:"requestIdentifier,omitempty"`
}

// MassExtendRenewalDateStatusResponse indicates the current status of a request to extend the subscription renewal date to all eligible subscribers.
// See https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldatestatusresponse
type MassExtendRenewalDateStatusResponse struct {
	RequestIdentifier string `json:"requestIdentifier,omitempty"`
	Complete          bool   `json:"complete,omitempty"`
	CompleteDate      int64  `json:"completeDate,omitempty"`
	SucceededCount    int    `json:"succeededCount,omitempty"`
	FailedCount       int    `json:"failedCount,omitempty"`
}

// SendAttemptItem contains the success or error information and the date the App Store server records when it attempts to send a server notification to your server.
// See https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
type SendAttemptItem struct {
	AttemptDate       int64             `json:"attemptDate,omitempty"`
	SendAttemptResult SendAttemptResult `json:"sendAttemptResult,omitempty"`
}

// NotificationHistoryResponseItem contains the App Store server notification history record, including the signed notification payload and the result of the server's first send attempt.
// See https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
type NotificationHistoryResponseItem struct {
	SignedPayload string                                                    `json:"signedPayload,omitempty"`
	SendAttempts  []SendAttemptItem                                         `json:"sendAttempts,omitempty"`
	Payload       *appstoreservernotifications.ResponseBodyV2DecodedPayload `json:"-"`
}

// NotificationHistoryResponse contains the App Store Server Notifications history for your app.
// See https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponse struct {
	PaginationToken     string                            `json:"paginationToken,omitempty"`
	HasMore             bool                              `json:"hasMore,omitempty"`
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory,omitempty"`
}

// SendTestNotificationResponse contains the test notification token.
// See https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken,omitempty"`
}

// CheckTestNotificationResponse contains the contents of the test notification sent by the App Store server and the result from your server.
// See https://developer.apple.com/documentation/appstoreserverapi/checktestnotificationresponse
type CheckTestNotificationResponse struct {
	SignedPayload string                                                    `json:"signedPayload,omitempty"`
	SendAttempts  []SendAttemptItem                                         `json:"sendAttempts,omitempty"`
	Payload       *appstoreservernotifications.ResponseBodyV2DecodedPayload `json:"-"`
}

// LastTransactionsItem contains the most recent App Store-signed transaction information and App Store-signed renewal information for an auto-renewable subscription.
// See https://developer.apple.com/documentation/appstoreserverapi/lasttransactionsitem
type LastTransactionsItem struct {
	Status                SubscriptionStatus            `json:"status,omitempty"`
	OriginalTransactionID string                        `json:"originalTransactionId,omitempty"`
	SignedTransactionInfo string                        `json:"signedTransactionInfo,omitempty"`
	SignedRenewalInfo     string                        `json:"signedRenewalInfo,omitempty"`
	RenewalPayload        *JWSRenewalInfoDecodedPayload `json:"-"`
	TransactionPayload    *JWSTransactionDecodedPayload `json:"-"`
}

// SubscriptionGroupIdentifierItem contains information for auto-renewable subscriptions, including signed transaction information and signed renewal information, for one subscription group.
// See https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifieritem
type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier,omitempty"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions,omitempty"`
}

// StatusResponse contains status information for all of a customer's auto-renewable subscriptions in your app.
// See https://developer.apple.com/documentation/appstoreserverapi/statusresponse
type StatusResponse struct {
	Environment Environment                       `json:"environment,omitempty"`
	BundleID    string                            `json:"bundleId,omitempty"`
	AppAppleID  int64                             `json:"appAppleId,omitempty"`
	Data        []SubscriptionGroupIdentifierItem `json:"data,omitempty"`
}
