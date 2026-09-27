package services

import (
	"warframe-checker/internal/models"
	"warframe-checker/internal/storage"
)

type UserService struct {
	db *storage.DB
}

func NewUserService(db *storage.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	result, err := s.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Username, user.PasswordHash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}
