package handlers

import (
	cartService "cart_service/internal/services/cart"
	cartpb "cart_service/proto/gen"
	"cart_service/utils"
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

type handler struct {
	cartpb.UnimplementedCartServiceServer
	service cartService.Service
}

func NewHandler(service cartService.Service) *handler {
	return &handler{
		service: service,
	}

}

func (h *handler) AddToCart(ctx context.Context, req *cartpb.AddToCartRequest) (*cartpb.CartStandardResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapError(errors.New("missing authentication metadata"))
	}
	emails := md.Get("x-user-email")
	resp, err := h.service.AddToCart(ctx, emails[0], req)
	if err != nil {
		return nil, utils.MapError(err)
	}
	return &cartpb.CartStandardResponse{
		Success:    true,
		Message:    "added to cart successful",
		StatusCode: 200,
		Result: &cartpb.CartStandardResponse_CartData{
			CartData: resp,
		},
	}, nil

}
