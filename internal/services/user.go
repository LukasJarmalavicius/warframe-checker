package services

import (
	"encoding/json"
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

func (s *UserService) UpdateInventory(userID int, inventory []models.InventoryItem) error {
	inventoryJSON, err := json.Marshal(inventory)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE users SET inventoryJSON = ? WHERE id = ?", inventoryJSON, userID)
	return err
}

func (s *UserService) GetInventory(userID int) ([]models.InventoryItem, error) {
	var inventoryJSON []byte
	err := s.db.QueryRow(
		"SELECT inventoryJSON FROM users WHERE id = ?", userID,
	).Scan(&inventoryJSON)
	if err != nil {
		return nil, err
	}

	items := []models.InventoryItem{}
	if len(inventoryJSON) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(inventoryJSON, &items); err != nil {
		return nil, err
	}
	return items, nil
}
