package utils

import (
	"cart_service/internal/domain"
	cartpb "cart_service/proto/gen"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func DomainCartToProto(cart *domain.Cart) *cartpb.CartResponse {
	pbItems := make([]*cartpb.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		pbItems = append(pbItems, &cartpb.CartItem{
			ProductId:   item.ProductID,
			Category:    item.Category,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			ImageUrl:    item.ImageURL,
			Subtotal:    item.Subtotal,
		})
	}
	return &cartpb.CartResponse{
		Email:      cart.Email,
		Items:      pbItems,
		TotalItems: cart.TotalItems,
		Subtotal:   cart.Subtotal,
		CreatedAt:  timestamppb.New(cart.CreatedAt),
		UpdatedAt:  timestamppb.New(cart.UpdatedAt),
	}
}
