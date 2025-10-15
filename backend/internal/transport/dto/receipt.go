package dto

// StoreResponse represents a store in the API response
type StoreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ItemResponse represents an item in the API response
type ItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	PricePaid   float64 `json:"price_paid"`
	Subtotal    float64 `json:"subtotal"` // quantity * price_paid
}

// ReceiptResponse represents a receipt in the API response with calculated fields (for detail view)
type ReceiptResponse struct {
	ID          string         `json:"id"`
	Store       StoreResponse  `json:"store"`
	BoughtDate  string         `json:"bought_date"` // ISO 8601: YYYY-MM-DD
	Items       []ItemResponse `json:"items"`
	Subtotal    float64        `json:"subtotal"`
	Discounts   float64        `json:"discounts"`
	TotalAmount float64        `json:"total_amount"`
}

// ReceiptListItem represents a receipt in list views (for left sidebar)
type ReceiptListItem struct {
	ID          string  `json:"id"`
	StoreName   string  `json:"store_name"`
	ItemCount   int     `json:"item_count"`
	BoughtDate  string  `json:"bought_date"` // ISO 8601: YYYY-MM-DD
	TotalAmount float64 `json:"total_amount"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Quantity  float64 `json:"quantity"`
	PricePaid float64 `json:"price_paid"`
}
