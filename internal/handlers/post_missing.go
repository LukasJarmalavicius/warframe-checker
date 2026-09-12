package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"warframe-checker/internal/models"
)

func (h *Handler) PostMissing(w http.ResponseWriter, r *http.Request) {
	log.Println("/postMissing")

	start := time.Now()
	var items []models.InventoryItem

	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	missingItems := h.inventoryService.GetMissing(items)

	if err := json.NewEncoder(w).Encode(missingItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("postMissing took %s", time.Since(start))
}
