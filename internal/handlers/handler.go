package handlers

import (
	"warframe-checker/internal/cache"
)

type Handler struct {
	API_URL  string
	WFCD_API string
	cache    *cache.Cache
}

func NewHandler(apiURL, wfcdAPI string, cache *cache.Cache) *Handler {
	return &Handler{
		API_URL:  apiURL,
		WFCD_API: wfcdAPI,
		cache:    cache,
	}
}
