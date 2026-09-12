package services

import (
	"log"
	"strings"
	"time"
	"warframe-checker/internal/cache"
	"warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

type InventoryService struct {
	cache     *cache.Cache
	client    *httpclient.Client
	marketAPI string
}

func NewInventoryService(cache *cache.Cache, client *httpclient.Client, marketAPI string) *InventoryService {
	return &InventoryService{
		cache:     cache,
		client:    client,
		marketAPI: marketAPI,
	}
}

func (s *InventoryService) GetCurrentPrimes() []models.TrimmedItem {
	start := time.Now()
	data := s.cache.All()
	items := make([]models.TrimmedItem, 0, len(data))

	for _, item := range data {
		if item.Vaulted || !item.IsPrime {
			continue
		}
		items = append(items, models.TrimmedItem{
			Name:     item.Name,
			Category: item.Category,
			Vaulted:  item.Vaulted,
		})
	}

	resurgenceItems := models.VaultTrader{}
	if err := s.client.FetchJson("https://api.warframestat.us/pc/vaultTrader", &resurgenceItems); err != nil {
		return items
	}
	schedule := resurgenceItems.Schedule[len(resurgenceItems.Schedule)-2]
	words := strings.Fields(schedule.Item)

	frame1 := words[3] + " " + words[5]
	frame2 := words[4] + " " + words[5]
	
	items = append(items, models.TrimmedItem{
		Name:     frame1,
		Category: "Resurgence Frame",
		Vaulted:  false,
	})
	items = append(items, models.TrimmedItem{
		Name:     frame2,
		Category: "Resurgence Frame",
		Vaulted:  false,
	})

	log.Printf("cache: InventoryService.GetCurrentPrimes took %s returned %d items", time.Since(start), len(items))
	return items
}

func (s *InventoryService) GetMissing(inventory []models.InventoryItem) []models.PartialSet {
	start := time.Now()
	data := s.cache.All()
	trimedInventory := make([]string, 0, len(inventory))
	for _, item := range inventory {
		trimedInventory = append(trimedInventory, strings.ToLower(strings.TrimSpace(item.Name)))
	}

	ownedSets := make(map[string]bool, len(trimedInventory))
	for _, item := range trimedInventory {
		ownedSets[item] = true
	}

	var result []models.PartialSet
	for _, item := range data {
		if !item.IsPrime && len(item.Components) > 0 {
			continue
		}

		var have, missingParts []string
		for _, component := range item.Components {

			if component.Type == "Resource" {
				continue
			}
			fullName := item.Name + " " + component.Name
			if !ownedSets[strings.ToLower(fullName)] {
				missingParts = append(missingParts, component.Name)
			} else {
				have = append(have, component.Name)
			}
		}
		if len(missingParts) > 0 && len(have) > 0 {
			result = append(result, models.PartialSet{
				SetName:      item.Name,
				MissingParts: missingParts,
				MissingCount: len(missingParts),
			})
		}
	}

	log.Printf("cache: GetMissingItems took %s", time.Since(start))
	return result
}
