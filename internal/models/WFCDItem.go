package models

type WFCDItem struct {
	UniqueName string              `json:"uniqueName"`
	Name       string              `json:"name"`
	Category   string              `json:"category"`
	IsPrime    bool                `json:"isPrime"`
	Vaulted    bool                `json:"vaulted"`
	Components []WFCDItemComponent `json:"components"`
}

type WFCDItemComponent struct {
	Name   string `json:"name"`
	Ducats int    `json:"ducats"`
}
