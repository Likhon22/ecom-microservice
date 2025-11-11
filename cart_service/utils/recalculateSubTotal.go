package utils

import "cart_service/internal/domain"

func RecalculateSubTotal(cart *domain.Cart) {
	cart.TotalItems = 0
	cart.Subtotal = 0

	for _, item := range cart.Items {
		cart.TotalItems += item.Quantity
		cart.Subtotal += item.Subtotal
	}
}
