package main

import (
	"log"
	"net/http"

	"github.com/darkwingdck/adora-coding-assessment/internal/api"
	"github.com/darkwingdck/adora-coding-assessment/internal/db"
)

func main() {
	database, err := db.Init("data/app.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("database initialized")

	api := api.NewService()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", api.Test())

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
