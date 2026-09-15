package handlers

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", h.GetItems)
	mux.HandleFunc("GET /orders", h.GetOrders)
	mux.HandleFunc("GET /currentPrimes", h.GetCurrentPrimes)
	mux.HandleFunc("POST /missing", h.PostMissing)

	mux.HandleFunc("POST /inventory", h.PostInventory)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html") })
	mux.Handle("GET /dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("web/dist"))))
	mux.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/style.css") })

	return mux
}
