package services

import (
	"strings"
	"testing"
	"warframe-checker/internal/cache"
	"warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

func newTestCache(items ...models.WFCDItem) *cache.Cache {
	cache := cache.NewCache()
	for _, item := range items {
		cache.Set(strings.ToLower(item.Name), item)
	}
	return cache
}

func newTestClient() *httpclient.Client {
	return httpclient.NewDefaultClient()
}

func TestInventoryService_GetUnvaulted(t *testing.T) {
	c := newTestCache(models.WFCDItem{
		Name:    "Rhino Prime",
		Vaulted: true,
		IsPrime: true,
	}, models.WFCDItem{
		Name:    "Caliban Prime",
		Vaulted: false,
		IsPrime: true,
	}, models.WFCDItem{
		Name:    "Nekros",
		Vaulted: false,
		IsPrime: false,
	})

	svc := NewInventoryService(c, newTestClient(), "")
	vaulted := svc.GetUnvaulted()
	if len(vaulted) != 1 {
		t.Fatalf("expected 1 item, got %d", len(vaulted))
	}
	if vaulted[0].Vaulted {
		t.Error("expected unvaulted item")
	}
	if vaulted[0].Name != "Caliban Prime" {
		t.Errorf("expected Caliban Prime, got %s", vaulted[0].Name)
	}
}
