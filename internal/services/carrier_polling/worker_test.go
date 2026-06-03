package carrierpolling

import (
	"context"
	"errors"
	"testing"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func Test_processUser(t *testing.T) {
	ctx := context.Background()

	t.Run("active carrier, inactive entitlement -> updates to active", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{ID: uuid.New(), UserID: "u_1", Active: false}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_1"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusActive}, nil)

		reason := dto.EntitlementReasonCarrierActive
		mockStore.EXPECT().
			UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
				UserID:          "u_1",
				Active:          true,
				Source:          dto.EntitlementSourceCarrier,
				Reason:          &reason,
				ExpiresAt:       ent.ExpiresAt,
				LastEventTimeMs: ent.LastEventTimeMs,
			}).Return(nil)

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})

	t.Run("active carrier, already active entitlement -> no update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{UserID: "u_2", Active: true}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_2"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusActive}, nil)
		// UpsertEntitlement must NOT be called

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})

	t.Run("inactive carrier, active entitlement -> updates to inactive", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{UserID: "u_3", Active: true}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_3"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusInactive}, nil)

		reason := dto.EntitlementReasonCarrierInactive
		mockStore.EXPECT().
			UpsertEntitlement(ctx, dto.UpsertEntitlementCmd{
				UserID:          "u_3",
				Active:          false,
				Source:          dto.EntitlementSourceCarrier,
				Reason:          &reason,
				ExpiresAt:       nil,
				LastEventTimeMs: ent.LastEventTimeMs,
			}).Return(nil)

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})

	t.Run("inactive carrier, already inactive entitlement -> no update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{UserID: "u_4", Active: false}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_4"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusInactive}, nil)
		// UpsertEntitlement must NOT be called

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})

	t.Run("api_error from carrier -> no update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{UserID: "u_5", Active: true}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_5"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusApiError}, nil)
		// UpsertEntitlement must NOT be called

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})

	t.Run("carrier service error -> no update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent := &store.Entitlement{UserID: "u_6", Active: true}

		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_6"}).
			Return(nil, errors.New("service unavailable"))
		// UpsertEntitlement must NOT be called

		NewWorker(mockStore, mockCarrier).processUser(ctx, ent)
	})
}

func Test_poll(t *testing.T) {
	ctx := context.Background()

	t.Run("GetCarrierEntitlements error -> no carrier calls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		mockStore.EXPECT().
			GetCarrierEntitlements(ctx).
			Return(nil, errors.New("db error"))
		// GetMobileCarrierUserStatus must NOT be called

		NewWorker(mockStore, mockCarrier).poll(ctx)
	})

	t.Run("no carrier entitlements -> no carrier calls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		mockStore.EXPECT().
			GetCarrierEntitlements(ctx).
			Return([]*store.Entitlement{}, nil)
		// GetMobileCarrierUserStatus must NOT be called

		NewWorker(mockStore, mockCarrier).poll(ctx)
	})

	t.Run("carrier entitlements returned -> processes each user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)
		mockCarrier := mocks.NewMockMobileCarrierService(ctrl)

		ent1 := &store.Entitlement{UserID: "u_1", Active: true}
		ent2 := &store.Entitlement{UserID: "u_2", Active: true}

		mockStore.EXPECT().
			GetCarrierEntitlements(ctx).
			Return([]*store.Entitlement{ent1, ent2}, nil)

		// both users must be checked against the carrier
		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_1"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusActive}, nil)
		mockCarrier.EXPECT().
			GetMobileCarrierUserStatus(ctx, dto.GetMobileCarrierUserStatusCmd{UserID: "u_2"}).
			Return(&dto.GetMobileCarrierUserStatusResult{Status: dto.MobileCarrierUserStatusActive}, nil)
		// both are already active, so UpsertEntitlement must NOT be called

		NewWorker(mockStore, mockCarrier).poll(ctx)
	})
}
