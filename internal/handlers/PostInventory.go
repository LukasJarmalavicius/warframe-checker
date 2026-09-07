package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	client "warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

type InventoryRequest struct {
	Items []models.InventoryItem `json:"items"`
}

type InventoryItemStatus struct {
	IsVaulted     bool
	Ducats        int
	CheapestPrice int                  `json:"cheapestPrice"`
	OtherPrices   []models.OrderFilter `json:"otherPrices"`
}

type InventoryItemMissing struct {
	Set          string   `json:"set"`
	Missing      []string `json:"missing"`
	MissingCount int      `json:"missingCount"`
}

type InventoryItemResponse struct {
	Name   string              `json:"name"`
	Status InventoryItemStatus `json:"status"`
}

type InventoryResponse struct {
	Items   []InventoryItemResponse `json:"items"`
	Missing []InventoryItemMissing  `json:"missing"`
}

func getVaulted(h *Handler, itemName string) (bool, int) {
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

	isWhole := partName == "" || strings.EqualFold(partName, "set")

	var ducatCount int
	if isWhole {
		for _, component := range item.Components {
			ducatCount += component.Ducats
		}
	} else {
		for _, component := range item.Components {
			if strings.EqualFold(component.Name, partName) {
				ducatCount = component.Ducats
				break
			}
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

func getStatus(h *Handler, req models.InventoryItem) InventoryItemStatus {
	var price []models.OrderFilter
	var isVaulted bool
	var ducats int
	var wg sync.WaitGroup

	wg.Go(func() {
		price = getPrices(h, req.Name)
	})
	wg.Go(func() {
		isVaulted, ducats = getVaulted(h, req.Name)
	})
	wg.Wait()

	if price == nil {
		return InventoryItemStatus{}
	}

	return InventoryItemStatus{
		IsVaulted:     isVaulted,
		Ducats:        ducats,
		CheapestPrice: price[0].Platinum,
		OtherPrices:   price[1:],
	}
}

func getMissing(h *Handler, req InventoryRequest, maxMissing int) []InventoryItemMissing {
	names := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		names = append(names, item.Name)
	}

	var missing []InventoryItemMissing

	sets := h.cache.AlmostCompleteSets(names, maxMissing)
	for _, set := range sets {
		missing = append(missing, InventoryItemMissing{
			Set:          set.SetName,
			Missing:      set.Missing,
			MissingCount: set.MissingCount,
		})
	}
	return missing
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

	missing := getMissing(h, req, maxMissing)

	results := make(chan InventoryItemResponse, len(req.Items))
	var wg sync.WaitGroup

	for _, item := range req.Items {
		wg.Go(func() {
			status := getStatus(h, item)
			results <- InventoryItemResponse{
				Name:   item.Name,
				Status: status,
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	responses := []InventoryItemResponse{}
	for result := range results {
		responses = append(responses, result)
	}

	log.Printf("/inventory took %dms\n", time.Since(start).Milliseconds())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(InventoryResponse{Items: responses, Missing: missing})
}
