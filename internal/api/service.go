package api

import (
	"net/http"

	"github.com/darkwingdck/adora-coding-assessment/store"
)

type Service interface {
	// StoreWebhook webhook for in-app store
	StoreWebhook() http.HandlerFunc
	// MarketplaceRevoke webhook for marketplace to revoke subscriptions
	MarketplaceRevoke() http.HandlerFunc
	// GetEntitlement get entitlement by userID
	GetEntitlement() http.HandlerFunc
	// MockCarrier mocks mobile carrier method to get person subscription status
	MockCarrier() http.HandlerFunc
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{
		store: store,
	}
}
