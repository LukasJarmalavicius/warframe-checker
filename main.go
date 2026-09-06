package main

import (
	"log"
	"net/http"
	"os"

	"warframe-checker/internal/config"
	"warframe-checker/internal/handlers"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	log.Println("Starting server on :8080")
	h := &handlers.Handler{API_URL: config.Get("API_URL"), WFCD_JSON: config.Get("WFCD_JSON")}
	router := handlers.NewRouter(h)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("URL: %s", h.API_URL)
	log.Printf("WFCD_JSON: %s", h.WFCD_JSON)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Server failed:", err)
		os.Exit(1)
	}

}
