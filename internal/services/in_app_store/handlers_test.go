package inappstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	inappstore "github.com/darkwingdck/adora-coding-assessment/internal/services/in_app_store"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testUserID = "u_42"

// newCmd builds a base UpdateUserEntitlementCmd for the given event type.
func newCmd(eventType dto.EventType) dto.UpdateUserEntitlementCmd {
	return dto.UpdateUserEntitlementCmd{
		EventID:     "evt_1",
		UserID:      testUserID,
		Type:        eventType,
		ProductID:   dto.ProductIDPremiumMonthly,
		EventTimeMs: 1000,
	}
}

// expectNewEvent sets up the two calls that every handler goes through before
// dispatching: lock-check + store event creation.
func expectNewEvent(mockStore *mocks.MockStore, ctx context.Context, existing *store.Entitlement) {
	mockStore.EXPECT().
		GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: testUserID, WithLock: true}).
		Return(existing, nil)
	mockStore.EXPECT().
		CreateStoreEvent(ctx, gomock.Any()).
		Return(true, nil)
}

// expectUpsert verifies that UpsertEntitlement is called with the expected
// active flag and reason, then returns nil.
func expectUpsert(t *testing.T, mockStore *mocks.MockStore, ctx context.Context, active bool, reason dto.EntitlementReason) {
	t.Helper()
	mockStore.EXPECT().
		UpsertEntitlement(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, cmd dto.UpsertEntitlementCmd) error {
			require.Equal(t, active, cmd.Active)
			require.NotNil(t, cmd.Reason)
			require.Equal(t, reason, *cmd.Reason)
			return nil
		})
}

// expectNotification sets up the GetEntitlementByUserID + CreateNotification
// calls made by scheduleExpirationNotification.
func expectNotification(mockStore *mocks.MockStore, ctx context.Context) {
	mockStore.EXPECT().
		GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: testUserID}).
		Return(&store.Entitlement{ID: uuid.New(), UserID: testUserID}, nil)
	mockStore.EXPECT().
		CreateNotification(ctx, gomock.Any()).
		Return(nil)
}

// --- INITIAL_PURCHASE ---

func Test_handleInitialPurchase(t *testing.T) {
	ctx := context.Background()

	t.Run("creates active entitlement and schedules notification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, nil)
		expectUpsert(t, mockStore, ctx, true, dto.EntitlementReasonInitialPurchase)
		expectNotification(mockStore, ctx)

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeInitialPurchase))
		require.NoError(t, err)
	})
}

// --- RENEWAL ---

func Test_handleRenewal(t *testing.T) {
	ctx := context.Background()

	t.Run("reactivates entitlement and schedules notification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, nil)
		expectUpsert(t, mockStore, ctx, true, dto.EntitlementReasonRenewal)
		expectNotification(mockStore, ctx)

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeRenewal))
		require.NoError(t, err)
	})
}

// --- CANCELLATION ---

func Test_handleCancellation(t *testing.T) {
	ctx := context.Background()

	t.Run("deactivates entitlement and preserves expiresAt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		existing := time.Now().Add(15 * 24 * time.Hour)
		ent := &store.Entitlement{UserID: testUserID, Active: true, ExpiresAt: &existing}

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, ent)
		// handleCancellation fetches current entitlement (without lock) to read expiresAt
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: testUserID}).
			Return(ent, nil)
		mockStore.EXPECT().
			UpsertEntitlement(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, cmd dto.UpsertEntitlementCmd) error {
				require.False(t, cmd.Active)
				require.Equal(t, dto.EntitlementReasonCancellation, *cmd.Reason)
				require.Equal(t, ent.ExpiresAt, cmd.ExpiresAt) // expiresAt preserved
				return nil
			})

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeCancellation))
		require.NoError(t, err)
	})

	t.Run("creates inactive entitlement when no prior entitlement exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, nil)
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: testUserID}).
			Return(nil, nil) // no existing entitlement
		mockStore.EXPECT().
			UpsertEntitlement(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, cmd dto.UpsertEntitlementCmd) error {
				require.False(t, cmd.Active)
				require.Nil(t, cmd.ExpiresAt)
				return nil
			})

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeCancellation))
		require.NoError(t, err)
	})
}

// --- BILLING_ISSUE ---

func Test_handleBillingIssue(t *testing.T) {
	ctx := context.Background()

	t.Run("deactivates entitlement and preserves expiresAt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		existing := time.Now().Add(10 * 24 * time.Hour)
		ent := &store.Entitlement{UserID: testUserID, Active: true, ExpiresAt: &existing}

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, ent)
		mockStore.EXPECT().
			GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{UserID: testUserID}).
			Return(ent, nil)
		mockStore.EXPECT().
			UpsertEntitlement(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, cmd dto.UpsertEntitlementCmd) error {
				require.False(t, cmd.Active)
				require.Equal(t, dto.EntitlementReasonBillingIssue, *cmd.Reason)
				require.Equal(t, ent.ExpiresAt, cmd.ExpiresAt)
				return nil
			})

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeBillingIssue))
		require.NoError(t, err)
	})
}

// --- EXPIRATION ---

func Test_handleExpiration(t *testing.T) {
	ctx := context.Background()

	t.Run("deactivates entitlement and clears expiresAt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, nil)
		mockStore.EXPECT().
			UpsertEntitlement(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, cmd dto.UpsertEntitlementCmd) error {
				require.False(t, cmd.Active)
				require.Equal(t, dto.EntitlementReasonExpiration, *cmd.Reason)
				require.Nil(t, cmd.ExpiresAt)
				return nil
			})

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeExpiration))
		require.NoError(t, err)
	})
}

// --- UN_CANCELLATION ---

func Test_handleUnCancellation(t *testing.T) {
	ctx := context.Background()

	t.Run("reactivates entitlement and schedules notification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		withTx(mockStore)
		expectNewEvent(mockStore, ctx, nil)
		expectUpsert(t, mockStore, ctx, true, dto.EntitlementReasonUnCancellation)
		expectNotification(mockStore, ctx)

		err := inappstore.NewService(mockStore).UpdateUserEntitlement(ctx, newCmd(dto.EventTypeUnCancellation))
		require.NoError(t, err)
	})
}
