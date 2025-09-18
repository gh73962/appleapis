package appstoreservernotifications

// NotificationType represents the type of App Store Server Notification v2
// See https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
type NotificationType string

const (
	// ConsumptionRequest indicates that the customer initiated a refund request for a consumable
	// in-app purchase or auto-renewable subscription, and the App Store is requesting consumption data
	TypeConsumptionRequest NotificationType = "CONSUMPTION_REQUEST"

	// TypeDidChangeRenewalPref indicates that the customer made a change to their subscription plan
	TypeDidChangeRenewalPref NotificationType = "DID_CHANGE_RENEWAL_PREF"

	// TypeDidChangeRenewalStatus indicates that the customer made a change to the subscription renewal status
	TypeDidChangeRenewalStatus NotificationType = "DID_CHANGE_RENEWAL_STATUS"

	// TypeDidFailToRenew indicates that the subscription failed to renew due to a billing issue
	TypeDidFailToRenew NotificationType = "DID_FAIL_TO_RENEW"

	// TypeDidRenew indicates that the subscription successfully renewed
	TypeDidRenew NotificationType = "DID_RENEW"

	// TypeExpired indicates that a subscription expired
	TypeExpired NotificationType = "EXPIRED"

	// TypeExternalPurchaseToken applies only to apps that use the External Purchase API
	TypeExternalPurchaseToken NotificationType = "EXTERNAL_PURCHASE_TOKEN"

	// TypeGracePeriodExpired indicates that the billing grace period has ended without renewing the subscription
	TypeGracePeriodExpired NotificationType = "GRACE_PERIOD_EXPIRED"

	// TypeMetadataUpdate indicates you used the Change Subscription Metadata endpoint
	TypeMetadataUpdate NotificationType = "METADATA_UPDATE"

	// TypeMigration indicates you used the Migrate a Subscription to Advanced Commerce API endpoint
	TypeMigration NotificationType = "MIGRATION"

	// TypeOfferRedeemed indicates that a customer with an active subscription redeemed a subscription offer
	TypeOfferRedeemed NotificationType = "OFFER_REDEEMED"

	// TypeOneTimeCharge indicates the customer purchased a consumable, non-consumable, or non-renewing subscription
	TypeOneTimeCharge NotificationType = "ONE_TIME_CHARGE"

	// TypePriceChange indicates that you called the Change Subscription Price endpoint
	TypePriceChange NotificationType = "PRICE_CHANGE"

	// TypePriceIncrease indicates that the system has informed the customer of an auto-renewable subscription price increase
	TypePriceIncrease NotificationType = "PRICE_INCREASE"

	// TypeRefund indicates that the App Store successfully refunded a transaction
	TypeRefund NotificationType = "REFUND"

	// TypeRefundDeclined indicates the App Store declined a refund request
	TypeRefundDeclined NotificationType = "REFUND_DECLINED"

	// TypeRefundReversed indicates the App Store reversed a previously granted refund due to a dispute
	TypeRefundReversed NotificationType = "REFUND_REVERSED"

	// TypeRenewalExtended indicates the App Store extended the subscription renewal date for a specific subscription
	TypeRenewalExtended NotificationType = "RENEWAL_EXTENDED"

	// TypeRenewalExtension indicates that the App Store is attempting to extend the subscription renewal date
	TypeRenewalExtension NotificationType = "RENEWAL_EXTENSION"

	// TypeRevoke indicates that an in-app purchase the customer was entitled to through Family Sharing is no longer available
	TypeRevoke NotificationType = "REVOKE"

	// TypeSubscribed indicates that the customer subscribed to an auto-renewable subscription
	TypeSubscribed NotificationType = "SUBSCRIBED"

	// TypeTest is sent when you request it by calling the Request a Test Notification endpoint
	TypeTest NotificationType = "TEST"
)

// Subtype represents the subtype of App Store Server Notification v2
// See https://developer.apple.com/documentation/appstoreservernotifications/subtype
type Subtype string

const (
	// SubtypeAccepted indicates customer consented to price increase or was notified
	SubtypeAccepted Subtype = "ACCEPTED"

	// SubtypeActiveTokenReminder indicates a monthly reminder that external purchase token is still active
	SubtypeActiveTokenReminder Subtype = "ACTIVE_TOKEN_REMINDER"

	// SubtypeAutoRenewDisabled indicates customer or App Store turned off auto-renewal
	SubtypeAutoRenewDisabled Subtype = "AUTO_RENEW_DISABLED"

	// SubtypeAutoRenewEnabled indicates customer enabled auto-renewal
	SubtypeAutoRenewEnabled Subtype = "AUTO_RENEW_ENABLED"

	// SubtypeBillingRecovery indicates expired subscription successfully renewed after failure
	SubtypeBillingRecovery Subtype = "BILLING_RECOVERY"

	// SubtypeBillingRetry indicates subscription expired because renewal failed during billing retry period
	SubtypeBillingRetry Subtype = "BILLING_RETRY"

	// SubtypeCreated indicates Apple created a custom link token for external purchases
	SubtypeCreated Subtype = "CREATED"

	// SubtypeDowngrade indicates customer downgraded or cross-graded subscription
	SubtypeDowngrade Subtype = "DOWNGRADE"

	// SubtypeFailure indicates subscription renewal date extension failed for individual subscription
	SubtypeFailure Subtype = "FAILURE"

	// SubtypeGracePeriod indicates subscription failed to renew due to billing issue during grace period
	SubtypeGracePeriod Subtype = "GRACE_PERIOD"

	// SubtypeInitialBuy indicates customer purchased subscription for first time or received through Family Sharing
	SubtypeInitialBuy Subtype = "INITIAL_BUY"

	// SubtypePending indicates customer was informed of price increase but hasn't accepted it
	SubtypePending Subtype = "PENDING"

	// SubtypePriceIncrease indicates subscription expired because customer didn't consent to price increase
	SubtypePriceIncrease Subtype = "PRICE_INCREASE"

	// SubtypeProductNotForSale indicates subscription expired because product wasn't available for purchase
	SubtypeProductNotForSale Subtype = "PRODUCT_NOT_FOR_SALE"

	// SubtypeResubscribe indicates customer resubscribed or received access through Family Sharing
	SubtypeResubscribe Subtype = "RESUBSCRIBE"

	// SubtypeSummary indicates App Store completed renewal date extension request for all eligible subscribers
	SubtypeSummary Subtype = "SUMMARY"

	// SubtypeUpgrade indicates customer upgraded or cross-graded subscription
	SubtypeUpgrade Subtype = "UPGRADE"

	// SubtypeUnreported indicates Apple created token but didn't receive a report
	SubtypeUnreported Subtype = "UNREPORTED"

	// SubtypeVoluntary indicates subscription expired after customer turned off auto-renewal
	SubtypeVoluntary Subtype = "VOLUNTARY"
)
