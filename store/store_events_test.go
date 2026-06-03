package store_test

import (
	"context"
	"testing"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/stretchr/testify/require"
)

func Test_CreateStoreEvent(t *testing.T) {
	requireDB(t)
	ctx := context.Background()

	t.Run("creates new event and returns true", func(t *testing.T) {
		userID := "test_store_event_new"
		cleanup(t, userID)

		created, err := testStore.CreateStoreEvent(ctx, dto.CreateStoreEventCmd{
			UserID:      userID,
			EventID:     "evt_test_new_001",
			Type:        dto.EventTypeInitialPurchase,
			ProductID:   dto.ProductIDPremiumMonthly,
			EventTimeMs: 1000,
		})
		require.NoError(t, err)
		require.True(t, created)
	})

	t.Run("returns false for duplicate event_id", func(t *testing.T) {
		userID := "test_store_event_dup"
		cleanup(t, userID)

		cmd := dto.CreateStoreEventCmd{
			UserID:      userID,
			EventID:     "evt_test_dup_001",
			Type:        dto.EventTypeInitialPurchase,
			ProductID:   dto.ProductIDPremiumMonthly,
			EventTimeMs: 1000,
		}

		created, err := testStore.CreateStoreEvent(ctx, cmd)
		require.NoError(t, err)
		require.True(t, created)

		created, err = testStore.CreateStoreEvent(ctx, cmd)
		require.NoError(t, err)
		require.False(t, created)
	})
}
