package models

type InventoryItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type PartialSet struct {
	SetName      string   `json:"setName"`
	Owned        []string `json:"owned,omitempty"`
	MissingParts []string `json:"missingParts"`
	MissingCount int      `json:"missingCount"`
}
