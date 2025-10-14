package models

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type Receipt struct {
	ID         string    `json:"id"`
	StoreName  string    `json:"store_name"`
	BoughtDate string    `json:"bought_date"` // ISO 8601 format: YYYY-MM-DD
	Items      []Item    `json:"items"`
	Discounts  float64   `json:"discounts"`
}
