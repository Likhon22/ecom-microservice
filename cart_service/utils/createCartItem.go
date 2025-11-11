package utils

import (
	"cart_service/internal/domain"
	cartpb "cart_service/proto/gen"
)

func CreateCartItem(product *cartpb.Product, quantity int32) domain.CartItem {
	imageURL := ""
	if len(product.ImageUrls) > 0 {
		imageURL = product.ImageUrls[0]
	}

	return domain.CartItem{
		ProductID:   product.ProductId,
		Category:    product.Category,
		ProductName: product.Name,
		Price:       product.Price,
		Quantity:    quantity,
		ImageURL:    imageURL,
		Subtotal:    product.Price * float64(quantity),
	}
}
