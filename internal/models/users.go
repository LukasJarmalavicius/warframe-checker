package models

type User struct {
	ID              int
	Username        string
	PasswordHash    string
	StoredInventory []InventoryItem `json:"inventoryJSON"`
}
