package domain

import "time"

type Cart struct {
	Email      string     `json:"email" redis:"email"` // Use email as identifier
	Items      []CartItem `json:"items" redis:"items"`
	TotalItems int32      `json:"total_items" redis:"total_items"`
	Subtotal   float64    `json:"subtotal" redis:"subtotal"`
	CreatedAt  time.Time  `json:"created_at" redis:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" redis:"updated_at"`
}

type CartItem struct {
	ProductID   string  `json:"product_id" redis:"product_id"`
	Category    string  `json:"category" redis:"category"`
	ProductName string  `json:"product_name" redis:"product_name"`
	Price       float64 `json:"price" redis:"price"`
	Quantity    int32   `json:"quantity" redis:"quantity"`
	ImageURL    string  `json:"image_url" redis:"image_url"`
	Subtotal    float64 `json:"subtotal" redis:"subtotal"`
}
