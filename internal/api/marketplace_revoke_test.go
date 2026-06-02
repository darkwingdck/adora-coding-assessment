package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkwingdck/adora-coding-assessment/internal/api"
	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
	"github.com/darkwingdck/adora-coding-assessment/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_MarketplaceRevoke(t *testing.T) {
	t.Run("invalid JSON -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/marketplace/revoke", nil)
		rec := httptest.NewRecorder()

		api.NewService(nil, nil, nil).MarketplaceRevoke()(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("store error -> 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		mockStore.EXPECT().
			RevokeMarketplaceEntitlements(gomock.Any(), dto.RevokeMarketplaceEntitlementsCmd{
				UserIDs: []string{"u_42", "u_91"},
			}).Return(errors.New("db error"))

		body := marshalJSON(t, api.MarketplaceRevokeRequest{UserIDs: []string{"u_42", "u_91"}})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/marketplace/revoke", body)
		rec := httptest.NewRecorder()

		api.NewService(mockStore, nil, nil).MarketplaceRevoke()(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success -> 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mocks.NewMockStore(ctrl)

		mockStore.EXPECT().
			RevokeMarketplaceEntitlements(gomock.Any(), dto.RevokeMarketplaceEntitlementsCmd{
				UserIDs: []string{"u_42"},
			}).Return(nil)

		body := marshalJSON(t, api.MarketplaceRevokeRequest{UserIDs: []string{"u_42"}})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/marketplace/revoke", body)
		rec := httptest.NewRecorder()

		api.NewService(mockStore, nil, nil).MarketplaceRevoke()(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
