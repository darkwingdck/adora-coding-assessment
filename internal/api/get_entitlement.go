package api

import "net/http"

type EntitlementResponse struct {
	Active        bool   `json:"active"`
	Source        string `json:"source"`
	ExpiresAt     string `json:"expiresAt"`
	LastChangedAt string `json:"lastChangedAt"`
	Reason        string `json:"reason"`
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
		_ = id
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}
}
