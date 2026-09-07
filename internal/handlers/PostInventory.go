package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	client "warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

type InventoryRequest struct {
	Items []models.InventoryItem `json:"items"`
}

type InventoryItemResponse struct {
	Name          string               `json:"name"`
	IsVaulted     bool                 `json:"isVaulted"`
	Ducats        int                  `json:"ducats"`
	CheapestPrice int                  `json:"cheapestPrice"`
	OtherPrices   []models.OrderFilter `json:"otherPrices"`
}

type InventoryResponse struct {
	Items []InventoryItemResponse `json:"items"`
}

func getVaultedStatus(h *Handler, itemName string) (bool, int) {
	lower := strings.ReplaceAll(strings.ToLower(itemName), "_", " ")
	if !strings.Contains(lower, "prime") {
		return false, 0
	}

	parentName := lower
	partName := lower
	words := strings.Fields(lower)
	if len(words) > 2 {
		partName = strings.Title(words[len(words)-1])
		parentName = strings.Join(words[:2], " ")
	}

	var item models.WFCDItem
	if data, ok := h.cache.Get(parentName); ok {
		item = data
	}

	var ducatCount int
	for _, component := range item.Components {
		if component.Name == partName {
			log.Println(component)
			ducatCount = component.Ducats
		} else {
			ducatCount += component.Ducats
		}
	}
	return item.Vaulted, ducatCount
}

func getPrices(h *Handler, itemName string) []models.OrderFilter {
	var items models.OrderResponse
	itemName = strings.ToLower(itemName)
	itemName = strings.ReplaceAll(itemName, " ", "_")

	if strings.HasSuffix(itemName, "prime") {
		itemName += "_set"
	}

	if err := client.FetchJson(h.API_URL+"orders/item/"+itemName+"/top", &items); err != nil {
		log.Println(err)
		return nil
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
			var ducats int

			innerWg := sync.WaitGroup{}
			innerWg.Go(func() {
				price = getPrices(h, item.Name)
			})
			innerWg.Go(func() {
				isVaulted, ducats = getVaultedStatus(h, item.Name)
			})

			innerWg.Wait()

			if price == nil {
				return
			}

			results <- InventoryItemResponse{
				Name:          item.Name,
				IsVaulted:     isVaulted,
				CheapestPrice: price[0].Platinum,
				Ducats:        ducats,
				OtherPrices:   price[1:],
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
