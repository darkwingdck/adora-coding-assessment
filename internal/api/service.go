package api

import (
	"net/http"

	mobilecarrier "github.com/darkwingdck/adora-coding-assessment/internal/services/mobile_carrier"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

type Service interface {
	// StoreWebhook webhook for in-app store
	StoreWebhook() http.HandlerFunc
	// MarketplaceRevoke webhook for marketplace to revoke subscriptions
	MarketplaceRevoke() http.HandlerFunc
	// GetEntitlement get entitlement by userID
	GetEntitlement() http.HandlerFunc
}

type service struct {
	store         store.Store
	mobileCarrier mobilecarrier.Service
}

func NewService(store store.Store, mobileCarrier mobilecarrier.Service) Service {
	return &service{
		store:         store,
		mobileCarrier: mobileCarrier,
	}
}
