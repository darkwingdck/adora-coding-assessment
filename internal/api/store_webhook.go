package api

import (
	"encoding/json"
	"net/http"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

type StoreWebhookRequest struct {
	EventID     string `json:"eventId"`
	UserID      string `json:"userId"`
	Type        string `json:"type"`
	EventTimeMs int64  `json:"eventTimeMs"`
	ProductID   string `json:"productId"`
}

// StoreWebhook godoc
//
//	@Summary	Ingest store webhook event
//	@Tags		webhooks
//	@Accept		json
//	@Produce	json
//	@Param		body	body		StoreWebhookRequest	true	"Store webhook event"
//	@Success	200		{string}	string				"ok"
//	@Failure	400		{string}	string				"bad request"
//	@Failure	500		{string}	string				"internal server error"
//	@Router		/webhooks/store [post]
func (s *service) StoreWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StoreWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		err := s.inAppStore.UpdateUserEntitlement(ctx, dto.UpdateUserEntitlementCmd{
			EventID:     req.EventID,
			UserID:      req.UserID,
			Type:        dto.EventType(req.Type),
			EventTimeMs: req.EventTimeMs,
			ProductID:   dto.ProductID(req.ProductID),
		})

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
