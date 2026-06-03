package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/api"
	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/darkwingdck/adora-coding-assessment/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_GetEntitlement(t *testing.T) {
	t.Run("store error -> 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		mockStore.EXPECT().
			GetEntitlementByUserID(gomock.Any(), dto.GetEntitlementByUserIDCmd{UserID: "u_42"}).
			Return(nil, errors.New("db error"))

		rec, req := newGetEntitlementRequest("u_42")
		api.NewService(mockStore, nil, nil).GetEntitlement()(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("entitlement not found -> 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		mockStore.EXPECT().
			GetEntitlementByUserID(gomock.Any(), dto.GetEntitlementByUserIDCmd{UserID: "u_42"}).
			Return(nil, nil)

		rec, req := newGetEntitlementRequest("u_42")
		api.NewService(mockStore, nil, nil).GetEntitlement()(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("entitlement found -> 200 with correct body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		reason := dto.EntitlementReasonRenewal
		ent := &store.Entitlement{
			ID:            uuid.New(),
			UserID:        "u_42",
			Active:        true,
			Source:        dto.EntitlementSourceStore,
			Reason:        &reason,
			ExpiresAt:     &expiresAt,
			LastChangedAt: time.Now(),
		}

		mockStore.EXPECT().
			GetEntitlementByUserID(gomock.Any(), dto.GetEntitlementByUserIDCmd{UserID: "u_42"}).
			Return(ent, nil)

		rec, req := newGetEntitlementRequest("u_42")
		api.NewService(mockStore, nil, nil).GetEntitlement()(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp api.EntitlementResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.True(t, resp.Active)
		require.Equal(t, "STORE", resp.Source)
		require.NotNil(t, resp.ExpiresAt)
		require.NotNil(t, resp.Reason)
		require.Equal(t, "RENEWAL", *resp.Reason)
	})

	t.Run("entitlement without expiresAt and reason -> 200, both nil in response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		ent := &store.Entitlement{
			ID:            uuid.New(),
			UserID:        "u_42",
			Active:        false,
			Source:        dto.EntitlementSourceNone,
			LastChangedAt: time.Now(),
		}

		mockStore.EXPECT().
			GetEntitlementByUserID(gomock.Any(), dto.GetEntitlementByUserIDCmd{UserID: "u_42"}).
			Return(ent, nil)

		rec, req := newGetEntitlementRequest("u_42")
		api.NewService(mockStore, nil, nil).GetEntitlement()(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp api.EntitlementResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.False(t, resp.Active)
		require.Nil(t, resp.ExpiresAt)
		require.Nil(t, resp.Reason)
	})
}

func newGetEntitlementRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/entitlement", nil)
	req.SetPathValue("id", userID)
	return httptest.NewRecorder(), req
}
