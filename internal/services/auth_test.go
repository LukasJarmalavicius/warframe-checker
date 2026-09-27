package services

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func newTestAuthService(t *testing.T) *AuthService {
	t.Helper()
	testDB := newTestDB(t)
	return NewAuthService(testDB, NewUserService(testDB))
}

func TestAuth(t *testing.T) {
	authService := newTestAuthService(t)
	err := authService.Register("test", "test")
	if err != nil {
		t.Fatal("failed to register:", err)
	}
	token, err := authService.Login("test", "test")
	if err != nil {
		t.Fatal("failed to login:", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
	user, err := authService.Authenticate(token)
	if err != nil {
		t.Fatal("failed to authenticate:", err)
	}
	if user.Username != "test" {
		t.Fatalf("expected username test, got %s", user.Username)
	}
}

func TestRegisterValidation(t *testing.T) {
	authService := newTestAuthService(t)
	tests := []struct {
		name     string
		username string
		password string
		want     error
	}{
		{"empty username", "", "test", ErrInvalidInput},
		{"blank username", "   ", "test", ErrInvalidInput},
		{"empty password", "test", "", ErrInvalidInput},
		{"password too long", "test", strings.Repeat("a", 73), ErrPasswordTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Register(tt.username, tt.password)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestRegisterDuplicate(t *testing.T) {
	authService := newTestAuthService(t)
	if err := authService.Register("test", "test"); err != nil {
		t.Fatal("failed to register:", err)
	}
	err := authService.Register("test", "other")
	if !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("expected %v, got %v", ErrUsernameTaken, err)
	}
}

func TestLoginFailures(t *testing.T) {
	authService := newTestAuthService(t)
	if err := authService.Register("test", "test"); err != nil {
		t.Fatal("failed to register:", err)
	}
	tests := []struct {
		name     string
		username string
		password string
		want     error
	}{
		{"wrong password", "test", "wrong", ErrInvalidCredentials},
		{"unknown user", "nobody", "test", ErrInvalidCredentials},
		{"empty input", "", "", ErrInvalidInput},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.Login(tt.username, tt.password)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestAuthenticateInvalidToken(t *testing.T) {
	authService := newTestAuthService(t)
	_, err := authService.Authenticate("not-a-real-token")
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("expected %v, got %v", ErrInvalidSession, err)
	}
}

func TestLogout(t *testing.T) {
	authService := newTestAuthService(t)
	if err := authService.Register("test", "test"); err != nil {
		t.Fatal("failed to register:", err)
	}
	token, err := authService.Login("test", "test")
	if err != nil {
		t.Fatal("failed to login:", err)
	}
	if err := authService.Logout(token); err != nil {
		t.Fatal("failed to logout:", err)
	}
	_, err = authService.Authenticate(token)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("expected %v after logout, got %v", ErrInvalidSession, err)
	}
}

func TestExpiredSession(t *testing.T) {
	authService := newTestAuthService(t)
	if err := authService.Register("test", "test"); err != nil {
		t.Fatal("failed to register:", err)
	}
	user, err := authService.users.GetUserByUsername("test")
	if err != nil {
		t.Fatal("failed to get user:", err)
	}
	token := "expired-token"
	_, err = authService.db.Exec(
		"INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		hashToken(token), user.ID, time.Now().Add(-time.Hour).Unix(),
	)
	if err != nil {
		t.Fatal("failed to insert session:", err)
	}
	_, err = authService.Authenticate(token)
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("expected %v for expired session, got %v", ErrInvalidSession, err)
	}
}

func TestCheckPassword(t *testing.T) {
	authService := NewAuthService(nil, nil)

	result, err := authService.HashPassword("test")
	if err != nil {
		t.Fatal("failed to hash password:", err)
	}
	err = authService.CheckPassword("test", result)
	if err != nil {
		t.Fatal("failed to check password:", err)
	}
}
