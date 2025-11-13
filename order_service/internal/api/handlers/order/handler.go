package orderHandler

import (
	orderService "order_service/internal/service/order"
	orderpb "order_service/proto/gen"
)

type Handler struct {
	orderpb.UnimplementedOrderServiceServer
	service orderService.Service
}

func NewHandler(service orderService.Service) *Handler {
	return &Handler{
		service: service,
	}

}
