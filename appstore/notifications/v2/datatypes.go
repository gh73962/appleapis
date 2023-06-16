package notifications

// NotificationType see https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
type NotificationType string

const (
	ConsumptionRequest     NotificationType = "CONSUMPTION_REQUEST"
	DidChangeRenewalPref   NotificationType = "DID_CHANGE_RENEWAL_PREF"
	DidChangeRenewalStatus NotificationType = "DID_CHANGE_RENEWAL_STATUS"
	DidFailToRenew         NotificationType = "DID_FAIL_TO_RENEW"
	DidRenew               NotificationType = "DID_RENEW"
	Expired                NotificationType = "EXPIRED"
	GracePeriodExpired     NotificationType = "GRACE_PERIOD_EXPIRED"
	OfferRedeemed          NotificationType = "OFFER_REDEEMED"
	PriceIncrease          NotificationType = "PRICE_INCREASE"
	Refund                 NotificationType = "REFUND"
	RefundDeclined         NotificationType = "REFUND_DECLINED"
	RefundReversed         NotificationType = "REFUND_REVERSED"
	RenewalExtended        NotificationType = "RENEWAL_EXTENDED"
	RenewalExtension       NotificationType = "RENEWAL_EXTENSION"
	Revoke                 NotificationType = "REVOKE"
	Subscribed             NotificationType = "SUBSCRIBED"
	Test                   NotificationType = "TEST"
)

// Subtype see https://developer.apple.com/documentation/appstoreservernotifications/subtype
type Subtype string

const (
	Accepted             Subtype = "ACCEPTED"
	AutoRenewDisabled    Subtype = "AUTO_RENEW_DISABLED"
	AutoRenewEnabled     Subtype = "AUTO_RENEW_ENABLED"
	BillingRecovery      Subtype = "BILLING_RECOVERY"
	BillingRetry         Subtype = "BILLING_RETRY"
	Downgrade            Subtype = "DOWNGRADE"
	Failure              Subtype = "FAILURE"
	GracePeriod          Subtype = "GRACE_PERIOD"
	InitialBuy           Subtype = "INITIAL_BUY"
	Pending              Subtype = "PENDING"
	SubtypePriceIncrease Subtype = "PRICE_INCREASE"
	ProductNotForSale    Subtype = "PRODUCT_NOT_FOR_SALE"
	Resubscribe          Subtype = "RESUBSCRIBE"
	Summary              Subtype = "SUMMARY"
	Upgrade              Subtype = "UPGRADE"
	Voluntary            Subtype = "VOLUNTARY"
)

// ResponseBodyV2 see https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2
type ResponseBodyV2 struct {
	SignedPayload string `json:"signedPayload"`
}

// ResponseBodyV2DecodedPayload see https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
type ResponseBodyV2DecodedPayload struct {
	NotificationType NotificationType `json:"notificationType,omitempty"`
	Subtype          Subtype          `json:"subtype,omitempty"`
	Data             data             `json:"data,omitempty"`
	Summary          summary          `json:"summary,omitempty"`
	Version          string           `json:"version,omitempty"`
	SignedDate       int64            `json:"signedDate,omitempty"`
	NotificationUUID string           `json:"notificationUUID,omitempty"`
}

type data struct {
	AppAppleID            int64  `json:"appAppleId,omitempty"`
	BundleID              string `json:"bundleId,omitempty"`
	BundleVersion         string `json:"bundleVersion,omitempty"`
	Environment           string `json:"environment,omitempty"`
	SignedRenewalInfo     string `json:"signedRenewalInfo,omitempty"`
	SignedTransactionInfo string `json:"signedTransactionInfo,omitempty"`
	Status                int    `json:"status,omitempty"`
}

type summary struct {
	AppAppleID             int64  `json:"appAppleId,omitempty"`
	BundleID               string `json:"bundleId,omitempty"`
	RequestIdentifier      string `json:"requestIdentifier,omitempty"`
	Environment            string `json:"environment,omitempty"`
	ProductID              string `json:"productId,omitempty"`
	StorefrontCountryCodes string `json:"storefrontCountryCodes,omitempty"`
	FailedCount            int    `json:"failedCount,omitempty"`
	SucceededCount         int    `json:"succeededCount,omitempty"`
}

// JWSDecodedHeader https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type JWSDecodedHeader struct {
	Alg string   `json:"alg,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

type JWSNotification struct {
	Header    JWSDecodedHeader
	Payload   ResponseBodyV2DecodedPayload
	Signature string
}
