package api

import "net/http"

type CarrierPlanResponse struct {
	Status string `json:"status"`
}

// MockCarrier godoc
//
//	@Summary	Mock carrier plan status endpoint
//	@Tags		mock
//	@Produce	json
//	@Param		userId	query		string				true	"User ID"
//	@Success	200		{object}	CarrierPlanResponse	"Carrier plan status"
//	@Router		/mock/carrier/plan [get]
func (s *service) MockCarrier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
