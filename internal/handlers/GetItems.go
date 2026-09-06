package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"warframe-checker/internal/models"
)

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	log.Println("/items")
	resp, err := http.Get(h.API_URL + "/items")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var items models.ItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		log.Println(err)
		return
	}

	log.Printf("/items: returned %d items\n", len(items.Data))

	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Println(err)
		return
	}

}
