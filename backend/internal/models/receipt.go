package models

type Item struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type Receipt struct {
	StoreName   string  `json:"store_name"`
	Items       []Item  `json:"items"`
	Discounts   float64 `json:"discounts,omitempty"`
}
