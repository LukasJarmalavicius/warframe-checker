package handlers

import (
	"log"
	"net/http"
)

func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("/logout")

	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if err := h.authService.Logout(cookie.Value); err != nil {
			log.Println("failed to delete session:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	clearSessionCookie(w, r)
	w.WriteHeader(http.StatusNoContent)
}
