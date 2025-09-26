package appstoreserver

// API endpoints constants
const (
	// Base URLs for different environments
	ProductionBaseURL   = "https://api.storekit.itunes.apple.com"
	SandboxBaseURL      = "https://api.storekit-sandbox.itunes.apple.com"
	LocalTestingBaseURL = "https://local-testing-base-url"
	AppleRootCAURL      = "https://www.apple.com/appleca/AppleIncRootCertificate.cer"
	AppleRootCAG2URL    = "https://www.apple.com/certificateauthority/AppleRootCA-G2.cer"
	AppleRootCAG3URL    = "https://www.apple.com/certificateauthority/AppleRootCA-G3.cer"
)

// Environment represents the server environment, either sandbox or production.
// See https://developer.apple.com/documentation/appstoreserverapi/environment
type Environment string

const (
	// EnvironmentSandbox represents the sandbox environment
	EnvironmentSandbox Environment = "Sandbox"
	// EnvironmentProduction represents the production environment
	EnvironmentProduction   Environment = "Production"
	EnvironmentLocalTesting Environment = "LocalTesting"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsValid checks if the environment value is valid
func (e Environment) IsValid() bool {
	switch e {
	case EnvironmentSandbox, EnvironmentProduction, EnvironmentLocalTesting:
		return true
	default:
		return false
	}
}

func (e Environment) BaseURL() string {
	switch e {
	case EnvironmentProduction:
		return ProductionBaseURL
	case EnvironmentSandbox:
		return SandboxBaseURL
	case EnvironmentLocalTesting:
		return LocalTestingBaseURL
	default:
		return ""
	}
}

// VerificationStatus represents the status of verification
type VerificationStatus int

const (
	// VerificationStatusOK indicates successful verification
	VerificationStatusOK VerificationStatus = iota
	// VerificationStatusFailure indicates verification failure
	VerificationStatusFailure
	// VerificationStatusInvalidAppIdentifier indicates invalid app identifier
	VerificationStatusInvalidAppIdentifier
	// VerificationStatusInvalidCertificate indicates invalid certificate
	VerificationStatusInvalidCertificate
	// VerificationStatusInvalidChainLength indicates invalid certificate chain length
	VerificationStatusInvalidChainLength
	// VerificationStatusInvalidChain indicates invalid certificate chain
	VerificationStatusInvalidChain
	// VerificationStatusInvalidEnvironment indicates invalid environment
	VerificationStatusInvalidEnvironment
)

// String returns the string representation of verification status
func (v VerificationStatus) String() string {
	switch v {
	case VerificationStatusOK:
		return "OK"
	case VerificationStatusFailure:
		return "VERIFICATION_FAILURE"
	case VerificationStatusInvalidAppIdentifier:
		return "INVALID_APP_IDENTIFIER"
	case VerificationStatusInvalidCertificate:
		return "INVALID_CERTIFICATE"
	case VerificationStatusInvalidChainLength:
		return "INVALID_CHAIN_LENGTH"
	case VerificationStatusInvalidChain:
		return "INVALID_CHAIN"
	case VerificationStatusInvalidEnvironment:
		return "INVALID_ENVIRONMENT"
	default:
		return "UNKNOWN"
	}
}

// AutoRenewStatus indicates the current renewal status for an auto-renewable subscription.
// See https://developer.apple.com/documentation/appstoreserverapi/autorenewstatus
type AutoRenewStatus int

const (
	AutoRenewStatusOff AutoRenewStatus = 0
	AutoRenewStatusOn  AutoRenewStatus = 1
)

// ExpirationIntent indicates the reason why a subscription expired.
// See https://developer.apple.com/documentation/appstoreserverapi/expirationintent
type ExpirationIntent int

const (
	ExpirationIntentCustomerCanceled             ExpirationIntent = 1 // The customer canceled their subscription
	ExpirationIntentBillingError                 ExpirationIntent = 2 // Billing error; for example, the customer's payment information is no longer valid
	ExpirationIntentDidNotConsentToPriceIncrease ExpirationIntent = 3 // The customer didn't consent to an auto-renewable subscription price increase that requires their consent
	ExpirationIntentProductNotAvailable          ExpirationIntent = 4 // The product wasn't available for purchase at the time of renewal
	ExpirationIntentOther                        ExpirationIntent = 5 // The subscription expired for some other reason
)

// InAppOwnershipType indicates whether the customer is the purchaser of the product.
// See https://developer.apple.com/documentation/appstoreserverapi/inappownershiptype
type InAppOwnershipType string

const (
	InAppOwnershipTypeFamilyShared InAppOwnershipType = "FAMILY_SHARED"
	InAppOwnershipTypePurchased    InAppOwnershipType = "PURCHASED"
)

// SubscriptionOfferType indicates whether the offer is an introductory offer or a promotional offer.
// See https://developer.apple.com/documentation/appstoreserverapi/offertype
type SubscriptionOfferType int

const (
	OfferTypeIntroductory          SubscriptionOfferType = 1
	OfferTypePromotional           SubscriptionOfferType = 2
	OfferTypeSubscriptionOfferCode SubscriptionOfferType = 3
	OfferTypeWinBack               SubscriptionOfferType = 4
)

// RevocationReason indicates the reason for a refunded transaction.
// See https://developer.apple.com/documentation/appstoreserverapi/revocationreason
type RevocationReason int

const (
	RevocationReasonOtherIssue RevocationReason = 0 // The App Store refunded the transaction on behalf of the customer for other reasons, for example, an accidental purchase
	RevocationReasonAppIssue   RevocationReason = 1 // The App Store refunded the transaction on behalf of the customer due to an actual or perceived issue within your app
)

// TransactionReason indicates the reason for the purchase transaction.
// See https://developer.apple.com/documentation/appstoreserverapi/transactionreason
type TransactionReason string

const (
	TransactionReasonPurchase TransactionReason = "PURCHASE"
	TransactionReasonRenewal  TransactionReason = "RENEWAL"
)

// ProductType indicates the type of the app account token.
// See https://developer.apple.com/documentation/appstoreserverapi/type
type ProductType string

const (
	TypeAutoRenewableSubscription ProductType = "Auto-Renewable Subscription"
	TypeNonConsumable             ProductType = "Non-Consumable"
	TypeConsumable                ProductType = "Consumable"
	TypeNonRenewingSubscription   ProductType = "Non-Renewing Subscription"
)

// AccountTenure indicates how long the customer's account has been active.
// See https://developer.apple.com/documentation/appstoreserverapi/accounttenure
type AccountTenure int

const (
	AccountTenure3Days   AccountTenure = 1 // 0-3 days
	AccountTenure10Days  AccountTenure = 2 // 3-10 days
	AccountTenure30Days  AccountTenure = 3 // 10-30 days
	AccountTenure90Days  AccountTenure = 4 // 30-90 days
	AccountTenure180Days AccountTenure = 5 // 90-180 days
	AccountTenure365Days AccountTenure = 6 // 180-365 days
	AccountTenureOver365 AccountTenure = 7 // Over 365 days
)

// ConsumptionStatus indicates the status of the in-app purchase.
// See https://developer.apple.com/documentation/appstoreserverapi/consumptionstatus
type ConsumptionStatus int

const (
	ConsumptionStatusNotConsumed       ConsumptionStatus = 1
	ConsumptionStatusPartiallyConsumed ConsumptionStatus = 2
	ConsumptionStatusFullyConsumed     ConsumptionStatus = 3
)

// DeliveryStatus indicates whether the app delivered the consumable in-app purchase.
// See https://developer.apple.com/documentation/appstoreserverapi/deliverystatus
type DeliveryStatus int

const (
	DeliveryStatusDelivered      DeliveryStatus = 0
	DeliveryStatusQualityIssue   DeliveryStatus = 1
	DeliveryStatusWrongItem      DeliveryStatus = 2
	DeliveryStatusServerOutage   DeliveryStatus = 3
	DeliveryStatusCurrencyChange DeliveryStatus = 4
	DeliveryStatusOtherReason    DeliveryStatus = 5
)

// LifetimeDollarsPurchased indicates the total amount, in USD, the customer has spent on in-app purchases.
// See https://developer.apple.com/documentation/appstoreserverapi/lifetimedollarspurchased
type LifetimeDollarsPurchased int

const (
	LifetimeDollarsPurchasedZero     LifetimeDollarsPurchased = 1
	LifetimeDollarsPurchasedUpTo50   LifetimeDollarsPurchased = 2
	LifetimeDollarsPurchasedUpTo100  LifetimeDollarsPurchased = 3
	LifetimeDollarsPurchasedUpTo500  LifetimeDollarsPurchased = 4
	LifetimeDollarsPurchasedUpTo1000 LifetimeDollarsPurchased = 5
	LifetimeDollarsPurchasedUpTo2000 LifetimeDollarsPurchased = 6
	LifetimeDollarsPurchasedOver2000 LifetimeDollarsPurchased = 7
)

// LifetimeDollarsRefunded indicates the total amount, in USD, the customer has received from refunded in-app purchases.
// See https://developer.apple.com/documentation/appstoreserverapi/lifetimedollarsrefunded
type LifetimeDollarsRefunded int

const (
	LifetimeDollarsRefundedZero     LifetimeDollarsRefunded = 1
	LifetimeDollarsRefundedUpTo50   LifetimeDollarsRefunded = 2
	LifetimeDollarsRefundedUpTo100  LifetimeDollarsRefunded = 3
	LifetimeDollarsRefundedUpTo500  LifetimeDollarsRefunded = 4
	LifetimeDollarsRefundedUpTo1000 LifetimeDollarsRefunded = 5
	LifetimeDollarsRefundedUpTo2000 LifetimeDollarsRefunded = 6
	LifetimeDollarsRefundedOver2000 LifetimeDollarsRefunded = 7
)

// Platform indicates the platform where the customer used your app.
// See https://developer.apple.com/documentation/appstoreserverapi/platform
type Platform int

const (
	PlatformApple    Platform = 1
	PlatformNonApple Platform = 2
)

// PlayTime indicates the amount of time the customer used the app.
// See https://developer.apple.com/documentation/appstoreserverapi/playtime
type PlayTime int

const (
	PlayTime5Min       PlayTime = 1 // 0-5 minutes
	PlayTime1Hour      PlayTime = 2 // 5-60 minutes
	PlayTime6Hours     PlayTime = 3 // 1-6 hours
	PlayTime1Day       PlayTime = 4 // 6-24 hours
	PlayTime4Days      PlayTime = 5 // 1-4 days
	PlayTime16Days     PlayTime = 6 // 4-16 days
	PlayTimeOver16Days PlayTime = 7 // Over 16 days
)

// UserStatus indicates the status of the customer's account.
// See https://developer.apple.com/documentation/appstoreserverapi/userstatus
type UserStatus int

const (
	UserStatusActive        UserStatus = 1
	UserStatusSuspended     UserStatus = 2
	UserStatusTerminated    UserStatus = 3
	UserStatusLimitedAccess UserStatus = 4
)

// ExtendReasonCode indicates the reason for extending a subscription renewal date.
// See https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
type ExtendReasonCode int

const (
	ExtendReasonCodeCustomerSatisfy ExtendReasonCode = 1
	ExtendReasonCodeOther           ExtendReasonCode = 2
	ExtendReasonCodeServiceIssue    ExtendReasonCode = 3
)

// OrderLookupStatus indicates whether the order ID in the request is valid for your app.
// See https://developer.apple.com/documentation/appstoreserverapi/orderlookupstatus
type OrderLookupStatus int

const (
	// OrderLookupStatusValid indicates the order ID is valid
	OrderLookupStatusValid OrderLookupStatus = 0
	// OrderLookupStatusInvalid indicates the order ID is invalid
	OrderLookupStatusInvalid OrderLookupStatus = 1
)

// SendAttemptResult represents the success or error information the App Store server records when it attempts to send an App Store server notification to your server.
// See https://developer.apple.com/documentation/appstoreserverapi/sendattemptresult
type SendAttemptResult string

const (
	SendAttemptResultSuccess                      SendAttemptResult = "SUCCESS"
	SendAttemptResultTimedOut                     SendAttemptResult = "TIMED_OUT"
	SendAttemptResultTLSIssue                     SendAttemptResult = "TLS_ISSUE"
	SendAttemptResultCircularRedirect             SendAttemptResult = "CIRCULAR_REDIRECT"
	SendAttemptResultNoResponse                   SendAttemptResult = "NO_RESPONSE"
	SendAttemptResultSocketIssue                  SendAttemptResult = "SOCKET_ISSUE"
	SendAttemptResultUnsupportedCharset           SendAttemptResult = "UNSUPPORTED_CHARSET"
	SendAttemptResultInvalidResponse              SendAttemptResult = "INVALID_RESPONSE"
	SendAttemptResultPrematureClose               SendAttemptResult = "PREMATURE_CLOSE"
	SendAttemptResultUnsuccessfulHTTPResponseCode SendAttemptResult = "UNSUCCESSFUL_HTTP_RESPONSE_CODE"
	SendAttemptResultOther                        SendAttemptResult = "OTHER"
)

// SubscriptionStatus represents the status of an auto-renewable subscription.
// See https://developer.apple.com/documentation/appstoreserverapi/status
type SubscriptionStatus int

const (
	StatusActive             SubscriptionStatus = 1
	StatusExpired            SubscriptionStatus = 2
	StatusBillingRetry       SubscriptionStatus = 3
	StatusBillingGracePeriod SubscriptionStatus = 4
	StatusRevoked            SubscriptionStatus = 5
)

// String returns the string representation of the Status
func (s SubscriptionStatus) String() string {
	switch s {
	case StatusActive:
		return "1"
	case StatusExpired:
		return "2"
	case StatusBillingRetry:
		return "3"
	case StatusBillingGracePeriod:
		return "4"
	case StatusRevoked:
		return "5"
	default:
		return ""
	}
}

// RefundPreference indicates your preferred outcome for the refund request.
// See https://developer.apple.com/documentation/appstoreserverapi/refundpreference
type RefundPreference int

const (
	RefundPreferenceGrant        RefundPreference = 1 // You prefer that Apple grants the refund
	RefundPreferenceDecline      RefundPreference = 2 // You prefer that Apple declines the refund
	RefundPreferenceNoPreference RefundPreference = 3 // You have no preference whether Apple grants or declines the refund
)

// OfferDiscountType indicates the payment mode for subscription offers on an auto-renewable subscription.
// See https://developer.apple.com/documentation/appstoreserverapi/offerdiscounttype
type OfferDiscountType string

const (
	// OfferDiscountTypeFreeTrial indicates a free trial payment mode
	OfferDiscountTypeFreeTrial OfferDiscountType = "FREE_TRIAL"
	// OfferDiscountTypePayAsYouGo indicates customers pay over single or multiple billing periods
	OfferDiscountTypePayAsYouGo OfferDiscountType = "PAY_AS_YOU_GO"
	// OfferDiscountTypePayUpFront indicates customers pay up front
	OfferDiscountTypePayUpFront OfferDiscountType = "PAY_UP_FRONT"
)
