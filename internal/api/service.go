package api

import "net/http"

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

type service struct{}

func NewService() Service {
	return &service{}
}
