package store

import (
	"context"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	// Users
	UpsertUser(ctx context.Context, cmd dto.UpsertUserCmd) error

	// Entitlements
	UpsertEntitlement(ctx context.Context, cmd dto.UpsertEntitlementCmd) error
	ListEntitlements(ctx context.Context, cmd dto.ListEntitlementsCmd) ([]*Entitlement, error)

	UpdateEntitlement(ctx context.Context, cmd dto.UpdateEntitlementCmd) error

	RevokeMarketplaceEntitlements(ctx context.Context, cmd dto.RevokeMarketplaceEntitlementsCmd) error

	// StoreEvents
	CreateStoreEvent(ctx context.Context, cmd dto.CreateStoreEventCmd) (bool, error)

	// Notifications
	CreateNotification(ctx context.Context, cmd dto.CreateNotificationCmd) error
	GetPendingNotifications(ctx context.Context) ([]*Notification, error)
	MarkNotificationSent(ctx context.Context, cmd dto.MarkNotificationSentCmd) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{
		db: db,
	}
}
