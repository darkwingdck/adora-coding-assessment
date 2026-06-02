// Package main adora-coding-assessment API
//
//	@title			Adora Coding Assessment API
//	@version		1.0
//	@description	Test API service
//	@host			localhost:8080
//	@BasePath		/
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/darkwingdck/adora-coding-assessment/config"
	_ "github.com/darkwingdck/adora-coding-assessment/docs"
	"github.com/darkwingdck/adora-coding-assessment/internal/api"
	carrierpolling "github.com/darkwingdck/adora-coding-assessment/internal/services/carrier_polling"
	inappstore "github.com/darkwingdck/adora-coding-assessment/internal/services/in_app_store"
	mobilecarrier "github.com/darkwingdck/adora-coding-assessment/internal/services/mobile_carrier"
	notificationworker "github.com/darkwingdck/adora-coding-assessment/internal/services/notification_worker"
	"github.com/darkwingdck/adora-coding-assessment/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	ctx := context.Background()

	pool, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer pool.Close()

	log.Println("database initialized")

	storage := store.NewStore(pool)
	mobileCarrierService := mobilecarrier.NewService()
	inAppStoreService := inappstore.NewService(storage)

	carrierWorker := carrierpolling.NewWorker(storage, mobileCarrierService)
	go carrierWorker.Run(ctx)

	notifWorker := notificationworker.NewWorker(storage)
	go notifWorker.Run(ctx)

	api := api.NewService(storage, mobileCarrierService, inAppStoreService)

	mux := http.NewServeMux()

	// ==[ ROUTES ]==
	mux.HandleFunc("POST /webhooks/store", api.StoreWebhook())
	mux.HandleFunc("POST /webhooks/marketplace/revoke", api.MarketplaceRevoke())
	mux.HandleFunc("GET /users/{id}/entitlement", api.GetEntitlement())

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
