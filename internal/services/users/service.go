package users

import (
	"context"
	"fmt"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

type Service interface {
	GenerateTestUsers(ctx context.Context) error
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{
		store: store,
	}
}

type testUserSpec struct {
	userID string
	source dto.EntitlementSource
	reason *dto.EntitlementReason
}

func (s *service) GenerateTestUsers(ctx context.Context) error {
	initialPurchase := dto.EntitlementReasonInitialPurchase
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	specs := make([]testUserSpec, 0, 30)
	for i := 1; i <= 10; i++ {
		specs = append(specs, testUserSpec{
			userID: fmt.Sprintf("test-store-%02d", i),
			source: dto.EntitlementSourceStore,
			reason: &initialPurchase,
		})
		specs = append(specs, testUserSpec{
			userID: fmt.Sprintf("test-carrier-%02d", i),
			source: dto.EntitlementSourceCarrier,
			reason: nil,
		})
		specs = append(specs, testUserSpec{
			userID: fmt.Sprintf("test-marketplace-%02d", i),
			source: dto.EntitlementSourceMarketplace,
			reason: nil,
		})
	}

	return s.store.WithTransaction(ctx, func(tx store.Store) error {
		for _, spec := range specs {
			if err := tx.UpsertUser(ctx, dto.UpsertUserCmd{UserID: spec.userID}); err != nil {
				return fmt.Errorf("tx.UpsertUser %s: %w", spec.userID, err)
			}
			if err := tx.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
				UserID:    spec.userID,
				Source:    spec.source,
				Reason:    spec.reason,
				ExpiresAt: &expiresAt,
			}); err != nil {
				return fmt.Errorf("tx.UpsertEntitlement %s: %w", spec.userID, err)
			}
		}
		return nil
	})
}
