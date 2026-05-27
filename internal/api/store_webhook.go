package api

import "net/http"

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
		w.WriteHeader(http.StatusOK)
	}
}
