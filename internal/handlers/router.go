package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", h.GetItems)
	mux.HandleFunc("GET /orders", h.GetOrders)
	mux.HandleFunc("GET /currentPrimes", h.GetCurrentPrimes)

	mux.HandleFunc("POST /inventory", h.PostInventory)

	return mux
}
