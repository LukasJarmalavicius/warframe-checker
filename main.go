package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"warframe-checker/internal/cache"
	"warframe-checker/internal/config"
	"warframe-checker/internal/handlers"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to open log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := config.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	log.Printf("--- Server starting (pid %d) ---\n", os.Getpid())
	log.Println("Starting server on :8080")
	cache := cache.NewCache()
	if err := cache.Load(config.Get("WFCD_JSON")); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}

	h := handlers.NewHandler(config.Get("API_URL"), config.Get("WFCD_API"), cache)
	router := handlers.NewRouter(h)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("URL: %s", h.API_URL)
	log.Printf("WFCD_API: %s", h.WFCD_API)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Server failed:", err)
		os.Exit(1)
	}

}
