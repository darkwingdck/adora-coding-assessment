package dto

type EventType string

const (
	EventTypeInitialPurchase EventType = "INITIAL_PURCHASE"
	EventTypeRenewal         EventType = "RENEWAL"
	EventTypeCancellation    EventType = "CANCELLATION"
	EventTypeBillingIssue    EventType = "BILLING_ISSUE"
	EventTypeExpiration      EventType = "EXPIRATION"
	EventTypeUnCancellation  EventType = "UN_CANCELLATION"
)

type ProductID string

const (
	ProductIDPremiumMonthly = "premium_monthly"
	ProductIDPremiumYearly  = "premium_yearly"
)

type MobileCarrierUserStatus string

const (
	MobileCarrierUserStatusActive   MobileCarrierUserStatus = "active"
	MobileCarrierUserStatusInactive MobileCarrierUserStatus = "inactive"
	MobileCarrierUserStatusApiError MobileCarrierUserStatus = "api_error"
)

type EntitlementSource string

const (
	EntitlementSourceStore       EntitlementSource = "STORE"
	EntitlementSourceCarrier     EntitlementSource = "CARRIER"
	EntitlementSourceMarketplace EntitlementSource = "MARKETPLACE"
	EntitlementSourceNone        EntitlementSource = "NONE"
)

type EntitlementReason string

const (
	EntitlementReasonInitialPurchase   EntitlementReason = "INITIAL_PURCHASE"
	EntitlementReasonRenewal           EntitlementReason = "RENEWAL"
	EntitlementReasonCancellation      EntitlementReason = "CANCELLATION"
	EntitlementReasonBillingIssue      EntitlementReason = "BILLING_ISSUE"
	EntitlementReasonExpiration        EntitlementReason = "EXPIRATION"
	EntitlementReasonUnCancellation    EntitlementReason = "UN_CANCELLATION"
	EntitlementReasonCarrierActive     EntitlementReason = "CARRIER_ACTIVE"
	EntitlementReasonCarrierInactive   EntitlementReason = "CARRIER_INACTIVE"
	EntitlementReasonMarketplaceRevoke EntitlementReason = "MARKETPLACE_REVOKE"
)
