package handlers

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"time"
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

type meResponse struct {
	Username string `json:"username"`
}
