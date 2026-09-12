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

func TestInventoryService_GetCurrentPrimes(t *testing.T) {
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
	vaulted := svc.GetCurrentPrimes()
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

func TestInventoryService_GetMissing(t *testing.T) {
	c := newTestCache(models.WFCDItem{
		Name: "Rhino Prime",
		Components: []models.WFCDItemComponent{
			{Name: "Blueprint"},
			{Name: "Chassis"},
			{Name: "Neuroptics"},
			{Name: "Resource", Type: "Resource"},
		},
	}, models.WFCDItem{
		Name: "Nekros",
		Components: []models.WFCDItemComponent{
			{Name: "Blueprint"},
		},
	})

	svc := NewInventoryService(c, newTestClient(), "")
	result := svc.GetMissing([]models.InventoryItem{
		{Name: "Rhino Prime Blueprint", Quantity: 1},
		{Name: "Nekros Blueprint", Quantity: 1},
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 partial set, got %d", len(result))
	}

	set := result[0]
	if set.SetName != "Rhino Prime" {
		t.Errorf("expected Rhino Prime, got %s", set.SetName)
	}
	if set.MissingCount != 2 {
		t.Errorf("expected 2 missing parts, got %d", set.MissingCount)
	}
	if len(set.MissingParts) != 2 {
		t.Fatalf("expected 2 missing parts, got %d", len(set.MissingParts))
	}
	for _, part := range set.MissingParts {
		if part == "Blueprint" {
			t.Error("blueprint should not be missing")
		}
	}
}

func TestInventoryService_GetMissing_NoInventory(t *testing.T) {
	c := newTestCache(models.WFCDItem{
		Name: "Rhino Prime",
		Components: []models.WFCDItemComponent{
			{Name: "Blueprint"},
			{Name: "Chassis"},
		},
	})

	svc := NewInventoryService(c, newTestClient(), "")
	result := svc.GetMissing(nil)

	if len(result) != 1 {
		t.Fatalf("expected 1 partial set, got %d", len(result))
	}
	if result[0].MissingCount != 2 {
		t.Errorf("expected 2 missing parts, got %d", result[0].MissingCount)
	}
}
