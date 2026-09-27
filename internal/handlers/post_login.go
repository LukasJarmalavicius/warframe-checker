package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"warframe-checker/internal/services"
)

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("/login")

	creds, err := decodeCredentials(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(creds.Username, creds.Password)
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, services.ErrInvalidCredentials):
		log.Printf("failed login for %q", creds.Username)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case err != nil:
		log.Printf("login failed for %q: %v", creds.Username, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, r, token)

	log.Printf("user %q logged in", creds.Username)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(meResponse{Username: creds.Username}); err != nil {
		log.Println("failed to encode login response:", err)
	}
}
