package domain

import "time"

type Product struct {
	ProductID   string  `dynamodbav:"product_id"`
	Name        string  `dynamodbav:"name"`
	Description string  `dynamodbav:"description,omitempty"`
	Category    string  `dynamodbav:"category"`
	Price       float64 `dynamodbav:"price"`

	ImageURLs     []string  `dynamodbav:"image_urls,omitempty"`
	Status        string    `dynamodbav:"status"`
	CreatedBy     string    `dynamodbav:"created_by"`
	IsFeatured    bool      `dynamodbav:"is_featured"`
	Tags          []string  `dynamodbav:"tags,omitempty"`
	AverageRating float64   `dynamodbav:"average_rating"`
	TotalReviews  int       `dynamodbav:"total_reviews"`
	CreatedAt     time.Time `dynamodbav:"created_at"`
	UpdatedAt     time.Time `dynamodbav:"updated_at"`
}
