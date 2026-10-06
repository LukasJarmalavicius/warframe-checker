package services

import (
	"net/http"
	"net/http/httptest"
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

func newFakeWFCDAPI(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pc/vaultTrader", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"schedule": [
			{"expiry": "2026-10-01T18:00:00.000Z", "item": "M P V Banshee Mirage Prime Dual Pack"},
			{"expiry": "2026-10-29T18:00:00.000Z", "item": "M P V Protea Ivara Prime Dual Pack"},
			{"expiry": "2026-11-26T18:00:00.000Z"}
		]}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func newCurrentPrimesTestCache() *cache.Cache {
	return newTestCache(models.WFCDItem{
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
}

func TestInventoryService_GetCurrentPrimes(t *testing.T) {
	srv := newFakeWFCDAPI(t)
	svc := NewInventoryService(newCurrentPrimesTestCache(), newTestClient(), srv.URL)

	got := svc.GetCurrentPrimes()

	want := []models.TrimmedItem{
		{Name: "Caliban Prime", Category: "", Vaulted: false},
		{Name: "Protea Prime", Category: "Resurgence Frame", Vaulted: false},
		{Name: "Ivara Prime", Category: "Resurgence Frame", Vaulted: false},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item %d = %+v, want %+v", i, got[i], want[i])
		}
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
		IsPrime: true,
	}, models.WFCDItem{
		Name: "Nekros",
		Components: []models.WFCDItemComponent{
			{Name: "Blueprint"},
		},
	})

	svc := NewInventoryService(c, newTestClient(), "")

	items := []models.InventoryItem{
		{Name: "Rhino Prime Blueprint", Quantity: 1},
		{Name: "Nekros Blueprint", Quantity: 1},
	}

	result := svc.GetMissing(items)

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
			t.Errorf("blueprint should not be missing")
		}
	}
}

func TestInventoryService_GetVaultedStatus(t *testing.T) {
	c := newTestCache(models.WFCDItem{
		Name:    "Rhino Prime",
		Vaulted: true,
		IsPrime: true,
		Components: []models.WFCDItemComponent{
			{Name: "Blueprint", Ducats: 100},
		},
	}, models.WFCDItem{
		Name:    "Caliban Prime",
		Vaulted: false,
		IsPrime: true,
		Components: []models.WFCDItemComponent{
			{Name: "Chassis", Ducats: 15},
		},
	}, models.WFCDItem{
		Name:    "Nekros",
		Vaulted: false,
		IsPrime: false,
		Components: []models.WFCDItemComponent{
			{Name: "Neuroptics", Ducats: 0},
		},
	})

	tests := []struct {
		name           string
		input          string
		expectedBool   bool
		expectedDucats int
		wantErr        bool
	}{
		{"Rhino", "Rhino prime blueprint", true, 100, false},
		{"Caliban", "caliban prime chassis", false, 15, false},
		{"Nekros", "nekros systems", false, 0, true},
	}

	svc := NewInventoryService(c, newTestClient(), "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ducats, err := svc.GetVaultedStatus(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetVaultedStatus(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if result != tt.expectedBool || ducats != tt.expectedDucats {
				t.Errorf("GetVaultedStatus(%s) = %v and %d ducats; want %v and %d ducats", tt.input, result, ducats, tt.expectedBool, tt.expectedDucats)
			}
		})
	}
}
