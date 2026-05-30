package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/jackc/pgx/v5"
)

func (s *store) UpsertEntitlement(ctx context.Context, cmd dto.UpsertEntitlementCmd) error {
	return nil
}

func (s *store) GetEntitlementByUserID(ctx context.Context, cmd dto.GetEntitlementByUserIDCmd) (*Entitlement, error) {
	if cmd.UserID == "" {
		return nil, fmt.Errorf("empty userID")
	}
	var entitlement Entitlement

	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, active, source, reason, expires_at, last_changed_at, last_event_time_ms
		FROM entitlements
		WHERE user_id = $1
	`, cmd.UserID).Scan(
		&entitlement.ID,
		&entitlement.UserID,
		&entitlement.Active,
		&entitlement.Source,
		&entitlement.Reason,
		&entitlement.ExpiresAt,
		&entitlement.LastChangedAt,
		&entitlement.LastEventTimeMs,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("s.pool.QueryRow: %w", err)
	}

	return &entitlement, nil

}

func (s *store) UpdateEntitlement(ctx context.Context, cmd dto.UpdateEntitlementCmd) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE entitlements SET
		active = $1,
		source = $2::entitlement_source,
		reason = $3::entitlement_reason,
		expires_at = $4,
		last_changed_at = NOW(),
		last_event_time_ms = $5
		WHERE user_id = $6
	`, cmd.Active, cmd.Source, cmd.Reason, cmd.ExpiresAt, cmd.LastEventTimeMs, cmd.UserID)

	if err != nil {
		return fmt.Errorf("s.pool.Exec: %w", err)
	}
	return nil
}

func (s *store) RevokeMarketplaceEntitlements(ctx context.Context, cmd dto.RevokeMarketplaceEntitlementsCmd) error {
	return nil
}
