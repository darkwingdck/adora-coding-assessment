package api

import (
	"net/http"

	inappstore "github.com/darkwingdck/adora-coding-assessment/internal/services/in_app_store"
	mobilecarrier "github.com/darkwingdck/adora-coding-assessment/internal/services/mobile_carrier"
	"github.com/darkwingdck/adora-coding-assessment/internal/services/users"
	"github.com/darkwingdck/adora-coding-assessment/store"
)

type Service interface {
	// StoreWebhook webhook for in-app store
	StoreWebhook() http.HandlerFunc
	// MarketplaceRevoke webhook for marketplace to revoke subscriptions
	MarketplaceRevoke() http.HandlerFunc
	// GetEntitlement get entitlement by userID
	GetEntitlement() http.HandlerFunc
	// GenerateTestUsers creates 30 test users (10 per source)
	GenerateTestUsers() http.HandlerFunc
}

type service struct {
	store         store.Store
	mobileCarrier mobilecarrier.Service
	inAppStore    inappstore.Service
	users         users.Service
}

func NewService(store store.Store, mobileCarrier mobilecarrier.Service, inAppStore inappstore.Service, users users.Service) Service {
	return &service{
		store:         store,
		mobileCarrier: mobileCarrier,
		inAppStore:    inAppStore,
		users:         users,
	}
}
