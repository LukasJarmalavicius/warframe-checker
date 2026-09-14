package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"warframe-checker/internal/models"
)

type InventoryRequest struct {
	Items []models.InventoryItem `json:"items"`
}

func (h *Handler) PostInventory(w http.ResponseWriter, r *http.Request) {
	log.Println("/inventory")
	var ogReq InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&ogReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := ogReq.Items

	if len(req) == 0 {
		http.Error(w, "items cannot be empty", http.StatusBadRequest)
		return
	}

	maxMissingStr := r.URL.Query().Get("maxMissing")
	maxMissing := 1
	if maxMissingStr != "" {
		parsed, err := strconv.Atoi(maxMissingStr)
		if err != nil {
			http.Error(w, "maxMissing must be a number", http.StatusBadRequest)
			return
		}
		maxMissing = parsed
	}
	log.Printf("maxmissing = %d", maxMissing)

	missing := h.inventoryService.GetMissing(req)

	w.Header().Set("Content-Type", "application/x-ndjson")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]any{"missing": missing})
	flusher.Flush()

	results := make(chan models.InventoryResponse, len(req))

	go h.inventoryService.PostInventory(results, req)

	for result := range results {
		encoder.Encode(map[string]any{"item": result})
		flusher.Flush()
	}

}
