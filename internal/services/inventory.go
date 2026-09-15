package services

import (
	"errors"
	"log"
	"strings"
	"sync"
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
	trimmedInventory := make([]string, 0, len(inventory))
	for _, item := range inventory {
		trimmedItem := strings.ToLower(item.Name)
		words := strings.Fields(trimmedItem)
		if len(words) > 2 && words[2] != "blueprint" {
			trimmedItem = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(item.Name), " blueprint", ""))
		}
		trimmedInventory = append(trimmedInventory, trimmedItem)
	}

	ownedSets := make(map[string]bool, len(trimmedInventory))
	for _, item := range trimmedInventory {
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

func (s *InventoryService) GetVaultedStatus(itemName string) (bool, int, error) {
	lower := strings.ReplaceAll(strings.ToLower(itemName), "_", " ")
	if !strings.Contains(lower, "prime") {
		return false, 0, errors.New("item is not prime")
	}

	parentName := lower
	partName := lower
	words := strings.Fields(lower)
	if len(words) > 2 {
		partName = strings.Title(words[len(words)-1])
		parentName = strings.Join(words[:2], " ")
	}
	var item models.WFCDItem
	if data, ok := s.cache.Get(parentName); ok {
		item = data
	}
	var ducats int
	for _, component := range item.Components {
		if component.Type == "Resource" {
			continue
		}
		if partName == component.Name {
			ducats = component.Ducats
		}
	}

	return item.Vaulted, ducats, nil
}

func (s *InventoryService) PostInventory(ch chan<- models.InventoryResponse, items []models.InventoryItem) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Go(func() {
			vaulted, ducats, err := s.GetVaultedStatus(item.Name)
			if err != nil {
				log.Printf("failed to get vaulted status for %s: %v", item.Name, err)
			}
			ch <- models.InventoryResponse{
				Name:    item.Name,
				Vaulted: vaulted,
				Ducats:  ducats,
			}
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}
