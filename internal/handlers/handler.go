package handlers

import (
	"warframe-checker/internal/cache"
	"warframe-checker/internal/services"
)

type Handler struct {
	API_URL          string
	WFCD_API         string
	cache            *cache.Cache
	inventoryService *services.InventoryService
}

func NewHandler(apiURL, wfcdAPI string, cache *cache.Cache, inventoryService *services.InventoryService) *Handler {
	return &Handler{
		API_URL:          apiURL,
		WFCD_API:         wfcdAPI,
		cache:            cache,
		inventoryService: inventoryService,
	}
}
