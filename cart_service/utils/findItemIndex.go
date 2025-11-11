package utils

import "cart_service/internal/domain"

func FindCartIndex(cartItems []domain.CartItem, productId string) int {
	for i, item := range cartItems {
		if item.ProductID == productId {
			return i
		}
	}
	return -1

}
