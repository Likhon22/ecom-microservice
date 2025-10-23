package authhandler

import (
	"auth_service/internal/services/auth"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
)

type handler struct {
	userpb.UnimplementedAuthServiceServer
	service auth.Service
}

func NewHandler(service auth.Service) *handler {
	return &handler{
		service: service,
	}
}
func (h *handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	result, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, utils.MapError(err)

	}
	return &userpb.LoginResponse{Message: result}, nil
}
