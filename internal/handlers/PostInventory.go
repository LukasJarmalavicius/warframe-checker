package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"warframe-checker/internal/models"
)

type InventoryRequest struct {
	Items []models.InventoryItem `json:"items"`
}

type InventoryItemResponse struct {
	Name      string               `json:"name"`
	IsVaulted bool                 `json:"isVaulted"`
	Prices    []models.OrderFilter `json:"prices"`
}

type InventoryResponse struct {
	Items []InventoryItemResponse `json:"items"`
}

func getVaultedStatus(h *Handler, itemName string) bool {
	lower := strings.ReplaceAll(strings.ToLower(itemName), "_", " ")
	if !strings.Contains(lower, "prime") {
		return false
	}
	parentName := lower
	words := strings.Fields(lower)
	if len(words) > 2 {
		parentName = strings.Join(words[:2], " ")
	}
	var item models.WFCDItem
	queryURL := h.WFCD_API + "items/" + url.PathEscape(parentName) + "?only=vaulted"
	if err := fetchJson(queryURL, &item); err != nil {
		log.Println(err)
		return false
	}
	return item.Vaulted
}

func getPrices(h *Handler, itemName string) []models.OrderFilter {
	var items models.OrderResponse
	itemName = strings.ToLower(itemName)
	itemName = strings.ReplaceAll(itemName, " ", "_")

	if err := fetchJson(h.API_URL+"orders/item/"+itemName+"/top", &items); err != nil {
		log.Println(err)
		return nil
	}

	filter := make([]models.OrderFilter, 0, len(items.Data.Sell))
	for _, item := range items.Data.Sell {
		filter = append(filter, models.OrderFilter{
			Platinum: item.Platinum,
			Quantity: item.Quantity,
		})
	}
	return filter
}

func (h *Handler) PostInventory(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("/inventory")
	var req InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "items cannot be empty", http.StatusBadRequest)
		return
	}

	results := make(chan InventoryItemResponse, len(req.Items))
	var wg sync.WaitGroup

	log.Printf("PostInventory called with %d items\n", len(req.Items))
	for _, item := range req.Items {
		wg.Go(func() {
			var price []models.OrderFilter
			var isVaulted bool

			innerWg := sync.WaitGroup{}
			innerWg.Go(func() {
				price = getPrices(h, item.Name)
			})
			innerWg.Go(func() {
				isVaulted = getVaultedStatus(h, item.Name)
			})
			innerWg.Wait()

			if price == nil {
				return
			}

			results <- InventoryItemResponse{
				Name:      item.Name,
				IsVaulted: isVaulted,
				Prices:    price,
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var responses []InventoryItemResponse
	for r := range results {
		responses = append(responses, r)
	}

	log.Printf("/inventory took %dms\n", time.Since(start).Milliseconds())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(InventoryResponse{Items: responses})
}
