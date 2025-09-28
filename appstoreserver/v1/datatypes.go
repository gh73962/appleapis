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
	AppAccountToken             string                `json:"appAccountToken,omitempty"`
	AppTransactionID            string                `json:"appTransactionId,omitempty"`
	BundleID                    string                `json:"bundleId"`
	Currency                    string                `json:"currency"`
	Environment                 Environment           `json:"environment"`
	ExpiresDate                 int64                 `json:"expiresDate,omitempty"`
	InAppOwnershipType          InAppOwnershipType    `json:"inAppOwnershipType"`
	IsUpgraded                  bool                  `json:"isUpgraded"`
	OfferDiscountType           OfferDiscountType     `json:"offerDiscountType,omitempty"`
	OfferIdentifier             string                `json:"offerIdentifier,omitempty"`
	OfferPeriod                 string                `json:"offerPeriod,omitempty"`
	OfferType                   SubscriptionOfferType `json:"offerType,omitempty"`
	OriginalPurchaseDate        int64                 `json:"originalPurchaseDate"`
	OriginalTransactionID       string                `json:"originalTransactionId"`
	Price                       int64                 `json:"price"`
	ProductID                   string                `json:"productId"`
	PurchaseDate                int64                 `json:"purchaseDate"`
	Quantity                    int                   `json:"quantity"`
	RevocationDate              int64                 `json:"revocationDate,omitempty"`
	RevocationReason            *RevocationReason     `json:"revocationReason,omitempty"`
	SignedDate                  int64                 `json:"signedDate"`
	Storefront                  string                `json:"storefront"`
	StorefrontID                string                `json:"storefrontId"`
	SubscriptionGroupIdentifier string                `json:"subscriptionGroupIdentifier"`
	TransactionID               string                `json:"transactionId"`
	TransactionReason           TransactionReason     `json:"transactionReason"`
	Type                        ProductType           `json:"type"`
	WebOrderLineItemID          string                `json:"webOrderLineItemId,omitempty"`
	AdvancedCommerceInfo        *AdvancedCommerceInfo `json:"advancedCommerceInfo,omitempty"`
}

// GetExpiresDate returns the expiration date as a time.Time
func (t *JWSTransactionDecodedPayload) GetExpiresDate() time.Time {
	return time.UnixMilli(t.ExpiresDate)
}

// GetOriginalPurchaseDate returns the original purchase date as a time.Time
func (t *JWSTransactionDecodedPayload) GetOriginalPurchaseDate() time.Time {
	return time.UnixMilli(t.OriginalPurchaseDate)
}

// GetPurchaseDate returns the purchase date as a time.Time
func (t *JWSTransactionDecodedPayload) GetPurchaseDate() time.Time {
	return time.UnixMilli(t.PurchaseDate)
}

// GetRevocationDate returns the revocation date as a time.Time
func (t *JWSTransactionDecodedPayload) GetRevocationDate() time.Time {
	return time.UnixMilli(t.RevocationDate)
}

// GetSignedDate returns the signed date as a time.Time
func (t *JWSTransactionDecodedPayload) GetSignedDate() time.Time {
	return time.UnixMilli(t.SignedDate)
}

// JWSRenewalInfoDecodedPayload contains subscription renewal information signed by the App Store.
// See https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	AppAccountToken             string                `json:"appAccountToken,omitempty"`
	AppTransactionID            string                `json:"appTransactionId"`
	AutoRenewProductID          string                `json:"autoRenewProductId"`
	AutoRenewStatus             AutoRenewStatus       `json:"autoRenewStatus"`
	Currency                    string                `json:"currency"`
	EligibleWinBackOfferIDs     []string              `json:"eligibleWinBackOfferIds,omitempty"`
	Environment                 Environment           `json:"environment"`
	ExpirationIntent            ExpirationIntent      `json:"expirationIntent"`
	GracePeriodExpiresDate      int64                 `json:"gracePeriodExpiresDate,omitempty"`
	IsInBillingRetryPeriod      bool                  `json:"isInBillingRetryPeriod"`
	OfferDiscountType           OfferDiscountType     `json:"offerDiscountType,omitempty"`
	OfferIdentifier             string                `json:"offerIdentifier,omitempty"`
	OfferPeriod                 string                `json:"offerPeriod,omitempty"`
	OfferType                   SubscriptionOfferType `json:"offerType,omitempty"`
	OriginalTransactionID       string                `json:"originalTransactionId"`
	PriceIncreaseStatus         int                   `json:"priceIncreaseStatus,omitempty"`
	ProductID                   string                `json:"productId"`
	RecentSubscriptionStartDate int64                 `json:"recentSubscriptionStartDate"`
	RenewalDate                 int64                 `json:"renewalDate"`
	RenewalPrice                int64                 `json:"renewalPrice"`
	SignedDate                  int64                 `json:"signedDate"`
	AdvancedCommerceInfo        *AdvancedCommerceInfo `json:"advancedCommerceInfo,omitempty"`
}

// AdvancedCommerceInfo Renewal information that is present only for Advanced Commerce SKUs.
// See https://developer.apple.com/documentation/appstoreserverapi/advancedcommercerenewalinfo
type AdvancedCommerceInfo struct {
	ConsistencyToken   string                        `json:"consistencyToken,omitempty"`
	Descriptors        *AdvancedCommerceDescriptors  `json:"descriptors,omitempty"`
	Items              []AdvancedCommerceRenewalItem `json:"items,omitempty"`
	Period             string                        `json:"period,omitempty"`
	RequestReferenceID string                        `json:"requestReferenceId,omitempty"`
	TaxCode            string                        `json:"taxCode,omitempty"`
}

// AdvancedCommerceDescriptors see https://developer.apple.com/documentation/appstoreserverapi/advancedcommercedescriptors
type AdvancedCommerceDescriptors struct {
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// AdvancedCommerceRenewalItem see https://developer.apple.com/documentation/appstoreserverapi/advancedcommercerenewalitem
type AdvancedCommerceRenewalItem struct {
	SKU         string                 `json:"SKU,omitempty"`
	Description string                 `json:"description,omitempty"`
	DisplayName string                 `json:"displayName,omitempty"`
	Offer       *AdvancedCommerceOffer `json:"offer,omitempty"`
	Price       int64                  `json:"price,omitempty"`
}

// AdvancedCommerceOffer see https://developer.apple.com/documentation/appstoreserverapi/advancedcommerceoffer
type AdvancedCommerceOffer struct {
	Period      string `json:"period,omitempty"`
	PeriodCount int    `json:"periodCount,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// GetGracePeriodExpiresDate returns the grace period expiration date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetGracePeriodExpiresDate() time.Time {
	return time.UnixMilli(r.GracePeriodExpiresDate)
}

// GetRecentSubscriptionStartDate returns the recent subscription start date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetRecentSubscriptionStartDate() time.Time {
	return time.UnixMilli(r.RecentSubscriptionStartDate)
}

// GetSignedDate returns the signed date as a time.Time
func (r *JWSRenewalInfoDecodedPayload) GetSignedDate() time.Time {
	return time.UnixMilli(r.SignedDate)
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
	Revision           string                          `json:"revision"`
	HasMore            bool                            `json:"hasMore"`
	BundleID           string                          `json:"bundleId"`
	AppAppleID         int                             `json:"appAppleId"`
	Environment        Environment                     `json:"environment"`
	SignedTransactions []string                        `json:"signedTransactions"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// OrderLookupResponse includes the order lookup status and an array of signed transactions for the in-app purchases in the order.
// See https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             OrderLookupStatus               `json:"status"`
	SignedTransactions []string                        `json:"signedTransactions"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// RefundHistoryResponse contains an array of signed JSON Web Signature (JWS) refunded transactions, and paging information.
// See https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundHistoryResponse struct {
	SignedTransactions []string                        `json:"signedTransactions"`
	Revision           string                          `json:"revision"`
	HasMore            bool                            `json:"hasMore"`
	Payloads           []*JWSTransactionDecodedPayload `json:"-"`
}

// ExtendRenewalDateResponse indicates whether an individual renewal-date extension succeeded, and related details.
// See https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldateresponse
type ExtendRenewalDateResponse struct {
	OriginalTransactionID string `json:"originalTransactionId"`
	WebOrderLineItemID    string `json:"webOrderLineItemId"`
	Success               bool   `json:"success"`
	EffectiveDate         int64  `json:"effectiveDate"`
}

// MassExtendRenewalDateResponse indicates the server successfully received the subscription-renewal-date extension request.
// See https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldateresponse
type MassExtendRenewalDateResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
}

// MassExtendRenewalDateStatusResponse indicates the current status of a request to extend the subscription renewal date to all eligible subscribers.
// See https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldatestatusresponse
type MassExtendRenewalDateStatusResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
	Complete          bool   `json:"complete"`
	CompleteDate      int64  `json:"completeDate"`
	SucceededCount    int    `json:"succeededCount"`
	FailedCount       int    `json:"failedCount"`
}

// SendAttemptItem contains the success or error information and the date the App Store server records when it attempts to send a server notification to your server.
// See https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
type SendAttemptItem struct {
	AttemptDate       int64             `json:"attemptDate"`
	SendAttemptResult SendAttemptResult `json:"sendAttemptResult"`
}

// NotificationHistoryResponseItem contains the App Store server notification history record, including the signed notification payload and the result of the server's first send attempt.
// See https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
type NotificationHistoryResponseItem struct {
	SignedPayload string                                      `json:"signedPayload"`
	SendAttempts  []SendAttemptItem                           `json:"sendAttempts"`
	Payload       *appstoreservernotifications.DecodedPayload `json:"-"`
}

// NotificationHistoryResponse contains the App Store Server Notifications history for your app.
// See https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponse struct {
	PaginationToken     string                            `json:"paginationToken"`
	HasMore             bool                              `json:"hasMore"`
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory"`
}

// SendTestNotificationResponse contains the test notification token.
// See https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}

// CheckTestNotificationResponse contains the contents of the test notification sent by the App Store server and the result from your server.
// See https://developer.apple.com/documentation/appstoreserverapi/checktestnotificationresponse
type CheckTestNotificationResponse struct {
	SignedPayload string                                      `json:"signedPayload"`
	SendAttempts  []SendAttemptItem                           `json:"sendAttempts"`
	Payload       *appstoreservernotifications.DecodedPayload `json:"-"`
}

// LastTransactionsItem contains the most recent App Store-signed transaction information and App Store-signed renewal information for an auto-renewable subscription.
// See https://developer.apple.com/documentation/appstoreserverapi/lasttransactionsitem
type LastTransactionsItem struct {
	Status                SubscriptionStatus            `json:"status"`
	OriginalTransactionID string                        `json:"originalTransactionId"`
	SignedTransactionInfo string                        `json:"signedTransactionInfo"`
	SignedRenewalInfo     string                        `json:"signedRenewalInfo"`
	RenewalPayload        *JWSRenewalInfoDecodedPayload `json:"-"`
	TransactionPayload    *JWSTransactionDecodedPayload `json:"-"`
}

// SubscriptionGroupIdentifierItem contains information for auto-renewable subscriptions, including signed transaction information and signed renewal information, for one subscription group.
// See https://developer.apple.com/documentation/appstoreserverapi/subscriptiongroupidentifieritem
type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions"`
}

// StatusResponse contains status information for all of a customer's auto-renewable subscriptions in your app.
// See https://developer.apple.com/documentation/appstoreserverapi/statusresponse
type StatusResponse struct {
	Environment Environment                       `json:"environment"`
	BundleID    string                            `json:"bundleId"`
	AppAppleID  int64                             `json:"appAppleId"`
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
}
