package main

import (
	"log"
	"net/http"
	"os"

	"warframe-checker/internal/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	log.Println("Starting server on :8080")
	cfg := config.Get("API_URL")

	log.Printf("URL: %s", cfg)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Server failed:", err)
		os.Exit(1)
	}

}
