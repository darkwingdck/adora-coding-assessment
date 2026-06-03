package inappstore_test

import (
	"context"
	"testing"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	inappstore "github.com/darkwingdck/adora-coding-assessment/internal/services/in_app_store"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// withTx makes WithTransaction invoke the callback so inner store calls can
// be verified on the same mock.
func withTx(mockStore *mocks.MockStore) {
	mockStore.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(store.Store) error) error {
			return fn(mockStore)
		})
}

func Test_UpdateUserEntitlement(t *testing.T) {
	ctx := context.Background()
	const userID = "u_42"

	t.Run("stale event -> skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		// existing entitlement has a newer timestamp
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID, WithLock: true}).
			Return(&store.Entitlement{UserID: userID, LastEventTimeMs: 2000}, nil)
		// CreateStoreEvent and UpsertEntitlement must NOT be called

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, dto.UpdateUserEntitlementCmd{
			EventID:     "evt_old",
			UserID:      userID,
			Type:        dto.EventTypeRenewal,
			ProductID:   dto.ProductIDPremiumMonthly,
			EventTimeMs: 1000, // older than 2000
		})
		require.NoError(t, err)
	})

	t.Run("duplicate event -> skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID, WithLock: true}).
			Return(nil, nil)
		mockStore.EXPECT().
			CreateStoreEvent(ctx, gomock.Any()).
			Return(false, nil) // already processed
		// UpsertEntitlement must NOT be called

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, dto.UpdateUserEntitlementCmd{
			EventID:     "evt_dup",
			UserID:      userID,
			Type:        dto.EventTypeInitialPurchase,
			ProductID:   dto.ProductIDPremiumMonthly,
			EventTimeMs: 1000,
		})
		require.NoError(t, err)
	})

	t.Run("unknown event type -> error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: userID, WithLock: true}).
			Return(nil, nil)
		mockStore.EXPECT().
			CreateStoreEvent(ctx, gomock.Any()).
			Return(true, nil)

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, dto.UpdateUserEntitlementCmd{
			EventID:     "evt_wat",
			UserID:      userID,
			Type:        "NOT_A_REAL_EVENT",
			EventTimeMs: 1000,
		})
		require.Error(t, err)
	})
}
