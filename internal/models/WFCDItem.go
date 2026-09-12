package models

import "time"

type TrimmedItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Vaulted  bool   `json:"vaulted"`
}

type WFCDItem struct {
	Name       string              `json:"name"`
	Category   string              `json:"category"`
	IsPrime    bool                `json:"isPrime"`
	Vaulted    bool                `json:"vaulted"`
	Components []WFCDItemComponent `json:"components"`
}

type WFCDItemComponent struct {
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Ducats int    `json:"ducats"`
}

type VaultTrader struct {
	Schedule []VaultTraderItem `json:"schedule"`
}

type VaultTraderItem struct {
	Expiry     time.Time `json:"expiry"`
	Item       string    `json:"item"`
	UniqueName string    `json:"uniqueName"`
}
