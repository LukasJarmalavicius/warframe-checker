package services

import (
	"log"
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

func (s *InventoryService) GetUnvaulted() []models.TrimmedItem {
	start := time.Now()
	data := s.cache.All()
	items := make([]models.TrimmedItem, 0, len(data))
	for _, item := range data {
		if item.Vaulted || !item.IsPrime {
			continue
		}
		items = append(items, models.TrimmedItem{
			UniqueName: item.UniqueName,
			Name:       item.Name,
			Category:   item.Category,
			Vaulted:    item.Vaulted,
		})
	}
	log.Printf("cache: Unvaulted took %s returned %d items", time.Since(start), len(items))
	return items
}
