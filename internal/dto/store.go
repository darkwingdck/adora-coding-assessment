package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==[ Users ]==

type UpsertUserCmd struct {
	UserID string
}

// ==[ Entitlements ]==

type UpsertEntitlementCmd struct {
	UserID string
}

type GetEntitlementByUserIDCmd struct {
	UserID   string
	WithLock bool
}

type UpdateEntitlementCmd struct {
	UserID          string
	Active          bool
	Source          EntitlementSource
	Reason          *EntitlementReason
	ExpiresAt       *time.Time
	LastEventTimeMs int64
}

type RevokeMarketplaceEntitlementsCmd struct {
	UserIDs []string
}

// ==[ Store Events ]==

type CreateStoreEventCmd struct {
	UserID      string
	EventID     string
	Type        EventType
	ProductID   ProductID
	EventTimeMs int64
}

// ==[ Notification ]==

type CreateNotificationCmd struct {
	UserID        string
	EntitlementID uuid.UUID
	Type          string
	ScheduledFor  time.Time
}

type MarkNotificationSentCmd struct {
	NotificationID uuid.UUID
}
