package models

import "time"

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type Receipt struct {
	ID        string    `json:"id"`
	StoreName string    `json:"store_name"`
	Items     []Item    `json:"items"`
	Discounts float64   `json:"discounts"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
