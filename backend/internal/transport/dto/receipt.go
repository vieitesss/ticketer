package dto

// ItemResponse represents an item in the API response
type ItemResponse struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

// ReceiptResponse represents a receipt in the API response with calculated fields
type ReceiptResponse struct {
	StoreName   string         `json:"store_name"`
	Items       []ItemResponse `json:"items"`
	Subtotal    float64        `json:"subtotal"`
	Discounts   float64        `json:"discounts"`
	TotalAmount float64        `json:"total_amount"`
}

// ReceiptListItem represents a receipt in list views (without items)
type ReceiptListItem struct {
	ID          string  `json:"id"`
	StoreName   string  `json:"store_name"`
	TotalAmount float64 `json:"total_amount"`
	Discounts   float64 `json:"discounts"`
}
