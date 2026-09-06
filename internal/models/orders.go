package models

type OrderFilter struct {
	Platinum int `json:"platinum"`
	Quantity int `json:"quantity"`
}

type OrderResponse struct {
	Data OrderData `json:"data"`
}

type OrderData struct {
	Sell []Order `json:"sell"`
	Buy  []Order `json:"buy"`
}

type Order struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Platinum int    `json:"platinum"`
	Quantity int    `json:"quantity"`
}
