package handlers

import (
	"io"
	"log"
	"net/http"
)

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	log.Println("GetItems")
	resp, err := http.Get(h.API_URL + "/items")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(body)
}
