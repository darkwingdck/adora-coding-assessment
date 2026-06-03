package store

import (
	"context"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Store interface {
	WithTransaction(ctx context.Context, fn func(Store) error) error

	// Entitlements
	UpsertEntitlement(ctx context.Context, cmd dto.UpsertEntitlementCmd) error
	GetEntitlementByUserID(ctx context.Context, cmd dto.GetEntitlementByUserIDCmd) (*Entitlement, error)
	GetCarrierEntitlements(ctx context.Context) ([]*Entitlement, error)
	RevokeMarketplaceEntitlements(ctx context.Context, cmd dto.RevokeMarketplaceEntitlementsCmd) error
	SeedTestEntitlements(ctx context.Context) error

	// StoreEvents
	CreateStoreEvent(ctx context.Context, cmd dto.CreateStoreEventCmd) (bool, error)

	// Notifications
	CreateNotification(ctx context.Context, cmd dto.CreateNotificationCmd) error
	GetPendingNotifications(ctx context.Context) ([]*Notification, error)
	MarkNotificationSent(ctx context.Context, cmd dto.MarkNotificationSentCmd) error
}

type store struct {
	rawPool *pgxpool.Pool
	pool    querier
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{
		rawPool: db,
		pool:    db,
	}
}

func (s *store) WithTransaction(ctx context.Context, fn func(Store) error) error {
	tx, err := s.rawPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(&store{rawPool: s.rawPool, pool: tx}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
