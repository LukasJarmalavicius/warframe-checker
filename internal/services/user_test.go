package services

import (
	"testing"
	"warframe-checker/internal/models"
)

func TestUserService(t *testing.T) {
	db := newTestDB(t)
	userService := NewUserService(db)
	user := &models.User{Username: "test", PasswordHash: "adsf"}
	err := userService.CreateUser(user)
	if err != nil {
		t.Fatal("failed to create user:", err)
	}
	check, err := userService.GetUser(user.ID)
	if err != nil {
		t.Fatal("failed to get user:", err)
	}
	if check.Username != user.Username {
		t.Fatal("user ID mismatch")
	}
	check, err = userService.GetUserByUsername("test")
	if err != nil {
		t.Fatal("failed to get user by username:", err)
	}
	if check.Username != user.Username {
		t.Fatal("user ID mismatch")
	}
}
