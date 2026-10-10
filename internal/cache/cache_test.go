package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoad(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/Warframes", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"name": "Baruuk Prime", "category": "Warframes", "isPrime": true}]`))
	})
	mux.HandleFunc("/Weapons", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[
			{"name": "Braton Prime", "category": "Primary", "isPrime": true},
		]`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cache := NewCache()
	if err := cache.Load(srv.URL + "/"); err != nil {
		t.Fatalf("failed to load cache: %v", err)
	}

	if _, ok := cache.Get("baruuk prime"); !ok {
		t.Fatalf("expected baruuk prime to be in cache")
	}
}
