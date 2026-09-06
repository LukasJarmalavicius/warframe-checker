package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", h.GetItems)
	mux.HandleFunc("/orders", h.GetPrice)

	return mux
}
