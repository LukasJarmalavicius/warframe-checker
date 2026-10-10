package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
	"time"

	"warframe-checker/internal/models"
	"warframe-checker/internal/services"
)

const sessionCookieName = "session"

// matches sessionTTL in services/auth.go
const sessionCookieMaxAge = 30 * 24 * time.Hour

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// decodeCredentials only accepts JSON bodies, which (together with SameSite
// cookies) stops cross-site HTML forms from posting to the auth endpoints.
func decodeCredentials(w http.ResponseWriter, r *http.Request) (credentials, error) {
	var creds credentials
	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		return creds, errors.New("content type must be application/json")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<10)
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return creds, errors.New("invalid request body")
	}
	return creds, nil
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(sessionCookieMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// currentUser returns the logged in user. If there is none it writes the
// error response itself and returns false.
func (h *Handler) currentUser(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return nil, false
	}

	user, err := h.authService.Authenticate(cookie.Value)
	if errors.Is(err, services.ErrInvalidSession) {
		clearSessionCookie(w, r)
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return nil, false
	}
	if err != nil {
		log.Println("failed to authenticate session:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return nil, false
	}
	return user, true
}

type meResponse struct {
	Username string `json:"username"`
}
