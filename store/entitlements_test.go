package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/stretchr/testify/require"
)

func Test_UpsertEntitlement(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("creates new entitlement", func(t *testing.T) {
		userID := "test_upsert_new"
		cleanup(t, userID)

		err := testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID:          userID,
			Source:          dto.EntitlementSourceStore,
			Reason:          ptr(dto.EntitlementReasonInitialPurchase),
			LastEventTimeMs: 1000,
		})
		require.NoError(t, err)

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
		require.Equal(t, userID, ent.UserID)
		require.True(t, ent.Active)
		require.Equal(t, dto.EntitlementSourceStore, ent.Source)
		require.Equal(t, int64(1000), ent.LastEventTimeMs)
	})

	t.Run("updates existing entitlement on conflict", func(t *testing.T) {
		userID := "test_upsert_update"
		cleanup(t, userID)

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID:          userID,
			Source:          dto.EntitlementSourceStore,
			Reason:          ptr(dto.EntitlementReasonInitialPurchase),
			LastEventTimeMs: 1000,
		}))

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID:          userID,
			Source:          dto.EntitlementSourceCarrier,
			Reason:          ptr(dto.EntitlementReasonCarrierActive),
			LastEventTimeMs: 2000,
		}))

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
		require.Equal(t, dto.EntitlementSourceCarrier, ent.Source)
		require.Equal(t, int64(2000), ent.LastEventTimeMs)
	})
}

func Test_GetEntitlementByUserID(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("returns error for empty userID", func(t *testing.T) {
		_, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: ""})
		require.Error(t, err)
	})

	t.Run("returns nil for non-existent user", func(t *testing.T) {
		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: "no_such_user_xyz"})
		require.NoError(t, err)
		require.Nil(t, ent)
	})

	t.Run("returns entitlement for existing user", func(t *testing.T) {
		userID := "test_get_ent"
		cleanup(t, userID)

		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID:          userID,
			Source:          dto.EntitlementSourceStore,
			Reason:          ptr(dto.EntitlementReasonRenewal),
			ExpiresAt:       &expiresAt,
			LastEventTimeMs: 5000,
		}))

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
		require.Equal(t, userID, ent.UserID)
		require.True(t, ent.Active)
		require.Equal(t, dto.EntitlementSourceStore, ent.Source)
		require.NotNil(t, ent.ExpiresAt)
	})
}

func Test_GetCarrierEntitlements(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	carrierUser := "test_carrier_ent_1"
	storeUser := "test_carrier_ent_2"
	cleanup(t, carrierUser, storeUser)

	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: carrierUser,
		Source: dto.EntitlementSourceCarrier,
		Reason: ptr(dto.EntitlementReasonCarrierActive),
	}))
	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: storeUser,
		Source: dto.EntitlementSourceStore,
		Reason: ptr(dto.EntitlementReasonInitialPurchase),
	}))

	entitlements, err := testStore.GetCarrierEntitlements(ctx)
	require.NoError(t, err)

	var carrierFound bool
	for _, e := range entitlements {
		if e.UserID == storeUser {
			t.Fatalf("STORE entitlement (user %s) should not appear in carrier results", storeUser)
		}
		if e.UserID == carrierUser {
			carrierFound = true
		}
	}
	require.True(t, carrierFound, "CARRIER entitlement should be in results")
}

func Test_UpdateEntitlement(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	userID := "test_update_ent"
	cleanup(t, userID)

	require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
		UserID: userID,
		Source: dto.EntitlementSourceStore,
		Reason: ptr(dto.EntitlementReasonInitialPurchase),
	}))

	err := testStore.UpdateEntitlement(ctx, dto.UpdateEntitlementCmd{
		UserID:          userID,
		Active:          false,
		Source:          dto.EntitlementSourceNone,
		Reason:          ptr(dto.EntitlementReasonCancellation),
		LastEventTimeMs: 9999,
	})
	require.NoError(t, err)

	ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
	require.NoError(t, err)
	require.NotNil(t, ent)
	require.False(t, ent.Active)
	require.Equal(t, dto.EntitlementSourceNone, ent.Source)
	require.Equal(t, dto.EntitlementReasonCancellation, *ent.Reason)
	require.Equal(t, int64(9999), ent.LastEventTimeMs)
}

func Test_RevokeMarketplaceEntitlements(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("revokes marketplace entitlements", func(t *testing.T) {
		userID := "test_revoke_mp"
		cleanup(t, userID)

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID: userID,
			Source: dto.EntitlementSourceMarketplace,
			Reason: ptr(dto.EntitlementReasonInitialPurchase),
		}))

		require.NoError(t, testStore.RevokeMarketplaceEntitlements(ctx, dto.RevokeMarketplaceEntitlementsCmd{
			UserIDs: []string{userID},
		}))

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
		require.False(t, ent.Active)
		require.Equal(t, dto.EntitlementSourceNone, ent.Source)
		require.Equal(t, dto.EntitlementReasonMarketplaceRevoke, *ent.Reason)
	})

	t.Run("does not revoke non-marketplace entitlements", func(t *testing.T) {
		userID := "test_revoke_nonmp"
		cleanup(t, userID)

		require.NoError(t, testStore.UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
			UserID: userID,
			Source: dto.EntitlementSourceStore,
			Reason: ptr(dto.EntitlementReasonInitialPurchase),
		}))

		require.NoError(t, testStore.RevokeMarketplaceEntitlements(ctx, dto.RevokeMarketplaceEntitlementsCmd{
			UserIDs: []string{userID},
		}))

		ent, err := testStore.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID})
		require.NoError(t, err)
		require.NotNil(t, ent)
		require.True(t, ent.Active)
		require.Equal(t, dto.EntitlementSourceStore, ent.Source)
	})
}
