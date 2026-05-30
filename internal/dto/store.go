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

type ListEntitlementsCmd struct {
	Filter struct {
		UserID *string
		Source *string
	}
}

type UpdateEntitlementCmd struct {
	UserID          string
	Active          bool
	Source          string
	Reason          *string
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
	Type        string
	ProductID   string
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
