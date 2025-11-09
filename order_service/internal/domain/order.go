package domain

import "time"

type Order struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	TotalAmount float64   `db:"total_amount" json:"total_amount"`
	Status      string    `db:"status" json:"status"`
	PaymentID   string    `db:"payment_id,omitempty" json:"payment_id,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
