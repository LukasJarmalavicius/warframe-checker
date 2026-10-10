package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := h.currentUser(w, r)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(meResponse{Username: user.Username}); err != nil {
		log.Println("failed to encode me response:", err)
	}
}
