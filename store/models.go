package store

import (
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/google/uuid"
)

type User struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type Entitlement struct {
	ID              uuid.UUID              `db:"id"`
	UserID          string                 `db:"user_id"`
	Active          bool                   `db:"active"`
	Source          dto.EntitlementSource  `db:"source"`
	Reason          *dto.EntitlementReason `db:"reason"`
	ExpiresAt       *time.Time             `db:"expires_at"`
	LastChangedAt   time.Time              `db:"last_changed_at"`
	LastEventTimeMs int64                  `db:"last_event_time_ms"`
}

type StoreEvent struct {
	ID          uuid.UUID     `db:"id"`
	UserID      string        `db:"user_id"`
	EventID     string        `db:"event_id"`
	Type        dto.EventType `db:"type"`
	ProductID   string        `db:"product_id"`
	EventTimeMs int64         `db:"event_time_ms"`
	ProcessedAt time.Time     `db:"processed_at"`
}

type Notification struct {
	ID            uuid.UUID  `db:"id"`
	UserID        string     `db:"user_id"`
	EntitlementID uuid.UUID  `db:"entitlement_id"`
	Type          string     `db:"type"`
	ScheduledFor  time.Time  `db:"scheduled_for"`
	SentAt        *time.Time `db:"sent_at"`
}
