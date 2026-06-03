package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SeedTestEntitlementsResponse struct {
	Store       []string `json:"store"`
	Carrier     []string `json:"carrier"`
	Marketplace []string `json:"marketplace"`
}

// SeedTestEntitlements godoc
//
//	@Summary	Seed 30 test entitlements (10 per source)
//	@Tags		test
//	@Produce	json
//	@Success	200	{object}	SeedTestEntitlementsResponse	"Created user IDs grouped by source"
//	@Failure	500	{string}	string							"internal server error"
//	@Router		/test/seed [post]
func (s *service) SeedTestEntitlements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.store.SeedTestEntitlements(r.Context()); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		resp := SeedTestEntitlementsResponse{
			Store:       make([]string, 10),
			Carrier:     make([]string, 10),
			Marketplace: make([]string, 10),
		}
		for i := 1; i <= 10; i++ {
			resp.Store[i-1] = fmt.Sprintf("user_store_test_%d", i)
			resp.Carrier[i-1] = fmt.Sprintf("user_carrier_test_%d", i)
			resp.Marketplace[i-1] = fmt.Sprintf("user_marketplace_test_%d", i)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(resp)
	}
}
