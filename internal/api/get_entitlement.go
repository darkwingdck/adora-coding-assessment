package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/darkwingdck/adora-coding-assessment/internal/dto"
)

type EntitlementResponse struct {
	Active        bool    `json:"active"`
	Source        string  `json:"source"`
	ExpiresAt     *string `json:"expiresAt"`
	LastChangedAt string  `json:"lastChangedAt"`
	Reason        *string `json:"reason"`
}

// GetEntitlement godoc
//
//	@Summary	Get current entitlement state for a user
//	@Tags		entitlements
//	@Produce	json
//	@Param		id	path		string				true	"User ID"
//	@Success	200	{object}	EntitlementResponse	"Current entitlement"
//	@Failure	404	{string}	string				"user not found"
//	@Failure	500	{string}	string				"internal server error"
//	@Router		/users/{id}/entitlement [get]
func (s *service) GetEntitlement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := r.Context()
		entitlement, err := s.store.GetEntitlementByUserID(ctx, dto.GetEntitlementByUserIDCmd{
			UserID: id,
		})
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		if entitlement == nil {
			http.Error(w, "Entitlement not found", http.StatusNotFound)
			return
		}

		resp := EntitlementResponse{
			Active:        entitlement.Active,
			Source:        string(entitlement.Source),
			LastChangedAt: entitlement.LastChangedAt.Format(time.RFC3339),
		}
		if entitlement.ExpiresAt != nil {
			expiresAt := entitlement.ExpiresAt.Format(time.RFC3339)
			resp.ExpiresAt = &expiresAt
		}
		if entitlement.Reason != nil {
			reason := string(*entitlement.Reason)
			resp.Reason = &reason
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(resp)
	}
}
