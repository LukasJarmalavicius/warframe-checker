package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", h.GetItems)
	return mux
}
