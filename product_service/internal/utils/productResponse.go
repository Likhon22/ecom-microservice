package utils

import (
	"product_service/internal/domain"
	productpb "product_service/proto/gen"
)

func ProductResponse(productData *domain.Product) *productpb.CreateProductResponse {
	return &productpb.CreateProductResponse{
		ProductId:   productData.ProductID,
		Name:        productData.Name,
		Description: productData.Description,
		Category:    productData.Category,
		Price:       productData.Price,
		ImageUrls:   productData.ImageURLs,
		Status:      productData.Status,
		IsFeatured:  productData.IsFeatured,
		Tags:        productData.Tags,
	}

}
