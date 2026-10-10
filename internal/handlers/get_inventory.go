package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) GetInventory(w http.ResponseWriter, r *http.Request) {
	user, ok := h.currentUser(w, r)
	if !ok {
		return
	}

	items, err := h.userService.GetInventory(user.ID)
	if err != nil {
		log.Println("failed to load inventory:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(InventoryRequest{Items: items}); err != nil {
		log.Println("failed to encode inventory response:", err)
	}
}
