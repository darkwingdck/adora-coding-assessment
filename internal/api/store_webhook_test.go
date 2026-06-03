package api_test

import (
	"bytes"
	"encoding/json"
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

func Test_StoreWebhook(t *testing.T) {
	t.Run("invalid JSON -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/store", bytes.NewBufferString("not json"))
		rec := httptest.NewRecorder()

		api.NewService(nil, nil, nil).StoreWebhook()(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("inAppStore error -> 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockInApp := mocks.NewMockInAppStoreService(ctrl)

		mockInApp.EXPECT().
			UpdateUserEntitlement(gomock.Any(), dto.UpdateUserEntitlementCmd{
				EventID:     "evt_1",
				UserID:      "u_42",
				Type:        dto.EventTypeInitialPurchase,
				EventTimeMs: 1000,
				ProductID:   dto.ProductID("premium_monthly"),
			}).Return(errors.New("service error"))

		body := marshalJSON(t, api.StoreWebhookRequest{
			EventID:     "evt_1",
			UserID:      "u_42",
			Type:        "INITIAL_PURCHASE",
			EventTimeMs: 1000,
			ProductID:   "premium_monthly",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/store", body)
		rec := httptest.NewRecorder()

		api.NewService(nil, nil, mockInApp).StoreWebhook()(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("success -> 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockInApp := mocks.NewMockInAppStoreService(ctrl)

		mockInApp.EXPECT().UpdateUserEntitlement(gomock.Any(), gomock.Any()).Return(nil)

		body := marshalJSON(t, api.StoreWebhookRequest{
			EventID:     "evt_2",
			UserID:      "u_42",
			Type:        "RENEWAL",
			EventTimeMs: 2000,
			ProductID:   "premium_monthly",
		})
		req := httptest.NewRequest(http.MethodPost, "/webhooks/store", body)
		rec := httptest.NewRecorder()

		api.NewService(nil, nil, mockInApp).StoreWebhook()(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func marshalJSON(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(b)
}
