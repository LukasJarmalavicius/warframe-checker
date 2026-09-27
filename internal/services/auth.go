package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"warframe-checker/internal/models"
	"warframe-checker/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionTTL = 30 * 24 * time.Hour
	// bcrypt ignores everything past 72 bytes
	maxPasswordBytes = 72
)

type AuthService struct {
	db    *storage.DB
	users *UserService
}

func NewAuthService(db *storage.DB, users *UserService) *AuthService {
	return &AuthService{
		db:    db,
		users: users,
	}
}

var (
	ErrInvalidInput       = errors.New("username and password are required")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrPasswordTooLong    = errors.New("password must be at most 72 bytes")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid or expired session")
)

func (s *AuthService) Register(username, password string) error {
	if strings.TrimSpace(username) == "" || password == "" {
		return ErrInvalidInput
	}
	if len(password) > maxPasswordBytes {
		return ErrPasswordTooLong
	}

	_, err := s.users.GetUserByUsername(username)
	if err == nil {
		return ErrUsernameTaken
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return err
	}
	users := models.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}
	err = s.users.CreateUser(&users)
	if err != nil {
		// a concurrent registration can slip past the lookup above
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameTaken
		}
		return err
	}
	return nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	if strings.TrimSpace(username) == "" || password == "" {
		return "", ErrInvalidInput
	}
	user, err := s.users.GetUserByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	err = s.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := s.NewSession(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) NewSession(userID int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	now := time.Now()
	if _, err := s.db.Exec("DELETE FROM sessions WHERE expires_at <= ?", now.Unix()); err != nil {
		return "", err
	}

	_, err := s.db.Exec(
		"INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		hashToken(token), userID, now.Add(sessionTTL).Unix(),
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Authenticate(token string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT u.id, u.username, u.password_hash
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = ? AND s.expires_at > ?`,
		hashToken(token), time.Now().Unix(),
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) Logout(token string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token_hash = ?", hashToken(token))
	return err
}
