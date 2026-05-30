package dto

type UpdateUserEntitlementCmd struct {
	EventID     string
	UserID      string
	Type        EventType
	EventTimeMs int64
	ProductID   ProductID
}

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
