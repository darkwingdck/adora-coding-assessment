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
	"github.com/darkwingdck/adora-coding-assessment/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	ctx := context.Background()

	database, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("database initialized")

	storage := store.NewStore(database)

	api := api.NewService(storage)

	mux := http.NewServeMux()

	// ==[ ROUTES ]==
	mux.HandleFunc("POST /webhooks/store", api.StoreWebhook())
	mux.HandleFunc("POST /webhooks/marketplace/revoke", api.MarketplaceRevoke())
	mux.HandleFunc("GET /users/{id}/entitlement", api.GetEntitlement())
	mux.HandleFunc("GET /mock/carrier/plan", api.MockCarrier())

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
