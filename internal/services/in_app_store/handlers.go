package inappstore

import (
	"context"
	"fmt"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

func (s *service) handleEvent(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	switch cmd.Type {
	case dto.EventTypeInitialPurchase:
		return s.handleInitialPurchase(ctx, tx, cmd)
	case dto.EventTypeRenewal:
		return s.handleRenewal(ctx, tx, cmd)
	case dto.EventTypeCancellation:
		return s.handleCancellation(ctx, tx, cmd)
	case dto.EventTypeBillingIssue:
		return s.handleBillingIssue(ctx, tx, cmd)
	case dto.EventTypeExpiration:
		return s.handleExparation(ctx, tx, cmd)
	case dto.EventTypeUnCancellation:
		return s.handleUnCancellation(ctx, tx, cmd)
	default:
		return fmt.Errorf("unknown event type: %s", cmd.Type)
	}
}

func (s *service) getExpiresAtFromProductID(productID dto.ProductID) (*time.Time, error) {
	var expiresAt time.Time
	switch productID {
	case dto.ProductIDPremiumMonthly:
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	case dto.ProductIDPremiumYearly:
		expiresAt = time.Now().Add(365 * 24 * time.Hour)
	default:
		return nil, fmt.Errorf("unknown product_id: %s", productID)
	}
	return &expiresAt, nil
}

// INITIAL_PURCHASE - creates entitlement with active = true, expires_at = NOW() + month/year
func (s *service) handleInitialPurchase(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	expiresAt, err := s.getExpiresAtFromProductID(cmd.ProductID)
	if err != nil {
		return fmt.Errorf("s.getExpiresAtFromProductID: %w", err)
	}

	reason := dto.EntitlementReasonInitialPurchase
	err = tx.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID:          cmd.UserID,
		Source:          dto.EntitlementSourceStore,
		Reason:          &reason,
		ExpiresAt:       expiresAt,
		LastEventTimeMs: cmd.EventTimeMs,
	})
	if err != nil {
		return fmt.Errorf("tx.UpsertEntitlement: %w", err)
	}
	return nil
}

// RENEWAL - set active = true, expires_at += month/year
func (s *service) handleRenewal(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	expiresAt, err := s.getExpiresAtFromProductID(cmd.ProductID)
	if err != nil {
		return fmt.Errorf("s.getExpiresAtFromProductID: %w", err)
	}

	reason := dto.EntitlementReasonRenewal
	err = tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          cmd.UserID,
		Active:          true,
		Source:          dto.EntitlementSourceStore,
		Reason:          &reason,
		ExpiresAt:       expiresAt, // TODO fix renewal status
		LastEventTimeMs: cmd.EventTimeMs,
	})
	if err != nil {
		return fmt.Errorf("tx.UpdateEntitlement: %w", err)
	}
	return nil
}

// CANCELLATION - set active = false, leave expires_at as it is (so we can do UN_CANCELLATION)
func (s *service) handleCancellation(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	current, err := tx.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: cmd.UserID})
	if err != nil {
		return fmt.Errorf("tx.GetEntitlementByUserID: %w", err)
	}

	var expiresAt *time.Time
	if current != nil {
		expiresAt = current.ExpiresAt
	}

	reason := dto.EntitlementReasonCancellation
	if err := tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          cmd.UserID,
		Active:          false,
		Source:          dto.EntitlementSourceStore,
		Reason:          &reason,
		ExpiresAt:       expiresAt,
		LastEventTimeMs: cmd.EventTimeMs,
	}); err != nil {
		return fmt.Errorf("tx.UpdateEntitlement: %w", err)
	}
	return nil
}

// BILLING_ISSUE - set active = false, leave expires_at as it is
func (s *service) handleBillingIssue(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	current, err := tx.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: cmd.UserID})
	if err != nil {
		return fmt.Errorf("tx.GetEntitlementByUserID: %w", err)
	}

	var expiresAt *time.Time
	if current != nil {
		expiresAt = current.ExpiresAt
	}

	reason := dto.EntitlementReasonBillingIssue
	if err := tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          cmd.UserID,
		Active:          false,
		Source:          dto.EntitlementSourceStore,
		Reason:          &reason,
		ExpiresAt:       expiresAt,
		LastEventTimeMs: cmd.EventTimeMs,
	}); err != nil {
		return fmt.Errorf("tx.UpdateEntitlement: %w", err)
	}
	return nil
}

// EXPIRATION - set active = false, expires_at = nil
func (s *service) handleExparation(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	reason := dto.EntitlementReasonExpiration
	err := tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          cmd.UserID,
		Active:          false,
		Source:          dto.EntitlementSourceStore,
		Reason:          &reason,
		ExpiresAt:       nil,
		LastEventTimeMs: cmd.EventTimeMs,
	})
	if err != nil {
		return fmt.Errorf("tx.UpdateEntitlement: %w", err)
	}
	return nil
}

// UN_CANCELLATION - set active = true, expires_at leave as it is
func (s *service) handleUnCancellation(ctx context.Context, tx store.Store, cmd dto.UpdateUserEntitlementCmd) error {
	expiresAt, err := s.getExpiresAtFromProductID(cmd.ProductID)
	if err != nil {
		return fmt.Errorf("s.getExpiresAtFromProductID: %w", err)
	}

	reason := dto.EntitlementReasonUnCancellation
	err = tx.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          cmd.UserID,
		Source:          dto.EntitlementSourceStore,
		Active:          true,
		Reason:          &reason,
		ExpiresAt:       expiresAt,
		LastEventTimeMs: cmd.EventTimeMs,
	})
	return nil
}
