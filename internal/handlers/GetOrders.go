package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"warframe-checker/internal/models"
)

type OrderFilter struct {
	Platinum int `json:"platinum"`
	Quantity int `json:"quantity"`
}

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	url := h.API_URL + "/orders/item/"
	if category != "" {
		url += category + "/top"
		log.Println(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var items models.OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		log.Println(err)
		return
	}

	filter := make([]OrderFilter, 0, len(items.Data.Sell))
	for _, item := range items.Data.Sell {
		filter = append(filter, OrderFilter{
			Platinum: item.Platinum,
			Quantity: item.Quantity,
		})
	}

	log.Printf("/price: returned %d orders\n", len(items.Data.Sell))

	if err := json.NewEncoder(w).Encode(filter); err != nil {
		log.Println(err)
		return
	}
}
