package models

type ItemsResponse struct {
	Data []Item `json:"data"`
}

type Item struct {
	ID      string   `json:"id"`
	Slug    string   `json:"slug"`
	GameRef string   `json:"gameRef"`
	Tags    []string `json:"tags"`
	Vaulted bool     `json:"vaulted"`
}
