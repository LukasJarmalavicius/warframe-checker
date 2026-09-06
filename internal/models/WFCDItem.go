package models

type WFCDItem struct {
	UniqueName string `json:"uniqueName"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	IsPrime    bool   `json:"isPrime"`
	Vaulted    bool   `json:"vaulted"`
}
