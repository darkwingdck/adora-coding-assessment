package api

import (
	"encoding/json"
	"net/http"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

type MarketplaceRevokeRequest struct {
	UserIDs []string `json:"userIds"`
}

// MarketplaceRevoke godoc
//
//	@Summary	Bulk revoke marketplace-granted access
//	@Tags		webhooks
//	@Accept		json
//	@Produce	json
//	@Param		body	body		MarketplaceRevokeRequest	true	"List of user IDs to revoke"
//	@Success	200		{string}	string						"ok"
//	@Failure	400		{string}	string						"bad request"
//	@Failure	500		{string}	string						"internal server error"
//	@Router		/webhooks/marketplace/revoke [post]
func (s *service) MarketplaceRevoke() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MarketplaceRevokeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		err := s.store.RevokeMarketplaceEntitlements(ctx, dto.RevokeMarketplaceEntitlementsCmd{
			UserIDs: req.UserIDs,
		})
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
