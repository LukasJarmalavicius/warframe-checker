package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Server failed:", err)
		os.Exit(1)
	}

}
