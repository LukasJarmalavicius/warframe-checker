package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"warframe-checker/internal/services"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.Authenticate(cookie.Value)
	if errors.Is(err, services.ErrInvalidSession) {
		clearSessionCookie(w, r)
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println("failed to authenticate session:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(meResponse{Username: user.Username}); err != nil {
		log.Println("failed to encode me response:", err)
	}
}
