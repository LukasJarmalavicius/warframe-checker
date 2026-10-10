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

func TestUserService_Inventory(t *testing.T) {
	db := newTestDB(t)
	userService := NewUserService(db)
	user := &models.User{Username: "test", PasswordHash: "adsf"}
	if err := userService.CreateUser(user); err != nil {
		t.Fatal("failed to create user:", err)
	}

	items, err := userService.GetInventory(user.ID)
	if err != nil {
		t.Fatal("failed to get empty inventory:", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty inventory, got %d items", len(items))
	}

	saved := []models.InventoryItem{
		{Name: "Volt Prime Chassis", Quantity: 2},
		{Name: "Braton Prime Barrel", Quantity: 1},
	}
	if err := userService.UpdateInventory(user.ID, saved); err != nil {
		t.Fatal("failed to update inventory:", err)
	}

	items, err = userService.GetInventory(user.ID)
	if err != nil {
		t.Fatal("failed to get inventory:", err)
	}
	if len(items) != len(saved) {
		t.Fatalf("expected %d items, got %d", len(saved), len(items))
	}
	for i := range saved {
		if items[i] != saved[i] {
			t.Fatalf("item %d mismatch: got %+v, want %+v", i, items[i], saved[i])
		}
	}
}
