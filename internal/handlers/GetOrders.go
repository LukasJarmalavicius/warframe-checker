package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	client "warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	url := h.API_URL + "/orders/item/"
	if category != "" {
		url += category + "/top"
		log.Println(url)
	}

	var items models.OrderResponse
	if err := client.FetchJson(url, &items); err != nil {
		log.Println(err)
		return
	}

	filter := make([]models.OrderFilter, 0, len(items.Data.Sell))
	for _, item := range items.Data.Sell {
		if item.User.Status != "ingame" {
			continue
		}
		filter = append(filter, models.OrderFilter{
			Platinum: item.Platinum,
			Quantity: item.Quantity,
		})
	}

	log.Printf("/orders: returned %d orders\n", len(filter))

	if err := json.NewEncoder(w).Encode(filter); err != nil {
		log.Println(err)
		return
	}
}
