package handlers

import (
	cartService "cart_service/internal/services/cart"
	cartpb "cart_service/proto/gen"
)

type handler struct {
	cartpb.UnimplementedCartServiceServer
	service cartService.Service
}

func NewHandler(service cartService.Service) *handler {
	return  &handler{
		service:service ,
	}

}
