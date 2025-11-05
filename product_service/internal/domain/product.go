package domain

import "time"

type Product struct {
	ProductID     string    `json:"product_id" dynamodbav:"ProductID"`
	Name          string    `json:"name" dynamodbav:"name"`
	Description   string    `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Category      string    `json:"category" dynamodbav:"Category"`
	Price         float64   `json:"price" dynamodbav:"price"`
	ImageURLs     []string  `json:"image_urls,omitempty" dynamodbav:"image_urls,omitempty"`
	Status        string    `json:"status" dynamodbav:"status"`
	CreatedBy     string    `json:"created_by" dynamodbav:"created_by"`
	IsFeatured    bool      `json:"is_featured" dynamodbav:"is_featured"`
	Tags          []string  `json:"tags,omitempty" dynamodbav:"tags,omitempty"`
	AverageRating float64   `json:"average_rating" dynamodbav:"average_rating"`
	TotalReviews  int       `json:"total_reviews" dynamodbav:"total_reviews"`
	CreatedAt     time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" dynamodbav:"updated_at"`
}
