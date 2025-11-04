package utils

import (
	"product_service/internal/domain"

	"time"
)

func BuildProductItem(p *domain.Product) map[string]interface{} {
	now := time.Now()

	item := map[string]interface{}{
		"PK":          GenerateProductPK(p.Category),
		"SK":          GenerateProductSK(p.ProductID),
		"product_id":  p.ProductID,
		"name":        p.Name,
		"category":    p.Category,
		"price":       p.Price,
		"status":      p.Status,
		"created_by":  p.CreatedBy,
		"is_featured": p.IsFeatured,
		"created_at":  now.Format(time.RFC3339),
		"updated_at":  now.Format(time.RFC3339),
	}

	// Optional fields: include only if provided
	if p.Description != "" {
		item["description"] = p.Description
	}
	if len(p.ImageURLs) > 0 {
		item["image_urls"] = p.ImageURLs
	}
	if len(p.Tags) > 0 {
		item["tags"] = p.Tags
	}
	if p.TotalReviews > 0 {
		item["total_reviews"] = p.TotalReviews
	}
	if p.AverageRating > 0 {
		item["average_rating"] = p.AverageRating
	}

	return item
}
