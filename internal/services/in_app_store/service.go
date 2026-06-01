package inappstore

import (
	"context"
	"fmt"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

type Service interface {
	UpdateUserEntitlement(ctx context.Context, cmd dto.UpdateUserEntitlementCmd) error
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) UpdateUserEntitlement(ctx context.Context, cmd dto.UpdateUserEntitlementCmd) error {
	return s.store.WithTransaction(ctx, func(tx store.Store) error {
		entitlement, err := tx.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{
			UserID:   cmd.UserID,
			WithLock: true, // prevents races with 2 webhook handlers
		})
		if err != nil {
			return fmt.Errorf("tx.GetEntitlementByUserID: %w", err)
		}

		var lastEventTimeMs int64
		if entitlement != nil {
			lastEventTimeMs = entitlement.LastEventTimeMs
		}

		// Old event
		if cmd.EventTimeMs < lastEventTimeMs {
			return nil
		}

		eventCreated, err := tx.CreateStoreEvent(ctx, dto.CreateStoreEventCmd{
			UserID:      cmd.UserID,
			EventID:     cmd.EventID,
			Type:        cmd.Type,
			ProductID:   cmd.ProductID,
			EventTimeMs: cmd.EventTimeMs,
		})
		if err != nil {
			return fmt.Errorf("tx.CreateStoreEvent: %w", err)
		}
		// Duplicate
		if !eventCreated {
			return nil
		}

		return s.handleEvent(ctx, tx, cmd)
	})
}
