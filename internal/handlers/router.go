package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", h.GetItems)
	mux.HandleFunc("GET /orders", h.GetOrders)
	mux.HandleFunc("GET /unvaulted", h.GetUnvaulted)

	mux.HandleFunc("POST /inventory", h.PostInventory)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html") })
	mux.Handle("GET /dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("web/dist"))))

	return mux
}
