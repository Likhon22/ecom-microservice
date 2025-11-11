package utils

import "cart_service/internal/domain"

func UpdateItemQuantity(item *domain.CartItem, additionalQuantity int32) {
	item.Quantity += additionalQuantity
	item.Subtotal = item.Price * float64(item.Quantity)
}
