package userhanlder

import (
	"auth_service/internal/services/auth"
	"auth_service/internal/types"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
)

type handler struct {
	userpb.UnimplementedUserServiceServer
	service auth.Service
}

func NewHandler(service auth.Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateCustomer(ctx context.Context, req *userpb.CreateCustomerRequest) (*userpb.CreateCustomerResponse, error) {

	result, err := h.service.CreateCustomer(ctx, types.CreateCustomerInput{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Phone:     req.GetPhone(),
		Password:  req.GetPassword(),
		Address:   req.GetAddress(),
		AvatarURL: req.GetAvatarUrl(),
	})
	if err != nil {
		return nil, utils.MapError(err)

	}
	return result, nil
}
func (h *handler) GetCustomerByEmail(ctx context.Context, req *userpb.GetCustomerByEmailRequest) (*userpb.CreateCustomerResponse, error) {

	result, err := h.service.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.MapError(err)

	}
	return result, nil
}
