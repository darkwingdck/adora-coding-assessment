package inappstore

import (
	"context"
	"fmt"
	"slices"
	"time"

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

var activeTypes = []dto.EventType{
	dto.EventTypeInitialPurchase,
	dto.EventTypeRenewal,
	dto.EventTypeUnCancellation,
}

var updateExpiresAtTypes = []dto.EventType{
	dto.EventTypeInitialPurchase,
	dto.EventTypeRenewal,
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
		// Old event
		if cmd.EventTimeMs < entitlement.LastEventTimeMs {
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

		var expiresAt *time.Time
		if slices.Contains(updateExpiresAtTypes, cmd.Type) {
			expiresAt, err = s.getExpiresAtFromProductID(cmd.ProductID)
			if err != nil {
				return fmt.Errorf("s.getExpiresAtFromProductID: %w", err)
			}
		}

		err = tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
			UserID:          cmd.UserID,
			Active:          slices.Contains(activeTypes, cmd.Type),
			Source:          dto.EntitlementSourceStore,
			Reason:          s.getReasonFromEventType(cmd.Type),
			ExpiresAt:       expiresAt,
			LastEventTimeMs: time.Now().Unix(),
		})
		if err != nil {
			return fmt.Errorf("tx.UpdateEntitlement: %w", err)
		}
		return nil
	})
}

func (s *service) getExpiresAtFromProductID(productID dto.ProductID) (*time.Time, error) {
	var expiresAt time.Time
	switch productID {
	case dto.ProductIDPremiumMonthly:
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	case dto.ProductIDPremiumYearly:
		expiresAt = time.Now().Add(365 * 24 * time.Hour)
	default:
		return nil, fmt.Errorf("unknown product_id: %v", productID)
	}
	return &expiresAt, nil
}

func (s *service) getReasonFromEventType(eventType dto.EventType) *dto.EntitlementReason {
	var reason dto.EntitlementReason
	switch eventType {
	case dto.EventTypeInitialPurchase:
		reason = dto.EntitlementReasonInitialPurchase
	case dto.EventTypeRenewal:
		reason = dto.EntitlementReasonRenewal
	case dto.EventTypeCancellation:
		reason = dto.EntitlementReasonCancellation
	case dto.EventTypeBillingIssue:
		reason = dto.EntitlementReasonBillingIssue
	case dto.EventTypeExpiration:
		reason = dto.EntitlementReasonExpiration
	case dto.EventTypeUnCancellation:
		reason = dto.EntitlementReasonCancellation
	default:
		return nil
	}
	return &reason
}
