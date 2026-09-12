package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func newTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	srv := httptest.NewServer(handler)
	c := NewClient(rate.NewLimiter(rate.Inf, 1), srv.Client())
	return c, srv
}

func TestFetchJson_Success(t *testing.T) {
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"name": "Braton Prime"})
	})
	defer srv.Close()

	var got map[string]string
	if err := c.FetchJson(srv.URL, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["name"] != "Braton Prime" {
		t.Errorf("got %v", got)
	}
}

func TestFetchJson_Failure(t *testing.T) {
	c, srv := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer srv.Close()

	var got map[string]string
	if err := c.FetchJson(srv.URL, &got); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
