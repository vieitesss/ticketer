package models

type Product struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	StoreID string `json:"store_id"`
}
