package api

import (
	"net/http"
)

// GenerateTestUsers godoc
//
//	@Summary	Generate 30 test users (10 per source: STORE, CARRIER, MARKETPLACE)
//	@Tags		dev
//	@Success	204
//	@Failure	500	{string}	string	"internal server error"
//	@Router		/dev/generate-test-users [post]
func (s *service) GenerateTestUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.users.GenerateTestUsers(r.Context()); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
