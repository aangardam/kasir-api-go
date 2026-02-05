package models

type Transaction struct {
	ID          int `json:"id"`
	TotalAmount int `json:"total_amount"`
	// CreatedAt   string               `json:"created_at"`
	Detail []TransactionDetails `json:"detail"`
}

type TransactionDetails struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	SubTotal      int    `json:"subtotal"`
}

type CheckoutRequest struct {
	Items []ChackoutItem `json:"items"`
}

type ChackoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
