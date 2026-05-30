package store

import (
	"context"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

func (s *store) UpsertEntitlement(ctx context.Context, cmd dto.UpsertEntitlementCmd) error {
	return nil
}

func (s *store) ListEntitlements(ctx context.Context, cmd dto.ListEntitlementsCmd) ([]*Entitlement, error) {
	return nil, nil
}

func (s *store) UpdateEntitlement(ctx context.Context, cmd dto.UpdateEntitlementCmd) error {
	return nil
}

func (s *store) RevokeMarketplaceEntitlements(ctx context.Context, cmd dto.RevokeMarketplaceEntitlementsCmd) error {
	return nil
}
