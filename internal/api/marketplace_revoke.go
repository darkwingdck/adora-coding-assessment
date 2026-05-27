package api

import "net/http"

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
		w.WriteHeader(http.StatusOK)
	}
}
