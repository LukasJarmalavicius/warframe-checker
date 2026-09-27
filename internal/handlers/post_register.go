package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"warframe-checker/internal/services"
)

func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("/register")

	creds, err := decodeCredentials(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.authService.Register(creds.Username, creds.Password)
	switch {
	case errors.Is(err, services.ErrInvalidInput),
		errors.Is(err, services.ErrPasswordTooLong),
		errors.Is(err, services.ErrUsernameTaken):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		log.Printf("register failed for %q: %v", creds.Username, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// log the new user straight in
	token, err := h.authService.Login(creds.Username, creds.Password)
	if err != nil {
		log.Printf("login after register failed for %q: %v", creds.Username, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, r, token)

	log.Printf("registered user %q", creds.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(meResponse{Username: creds.Username}); err != nil {
		log.Println("failed to encode register response:", err)
	}
}
