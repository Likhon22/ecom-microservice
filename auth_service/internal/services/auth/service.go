package auth

import (
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/types"
	userpb "auth_service/proto/gen"
	"context"
	"fmt"
)

type Service interface {
	CreateCustomer(ctx context.Context, payload types.CreateCustomerInput) (*userpb.CreateCustomerResponse, error)
	// GetCustomerByEmail(ctx context.Context, email string) (*types.CreateCustomerResult, error)
	// GetCustomers(ctx context.Context) ([]*types.CreateCustomerResult, error)
	// DeleteCustomer(ctx context.Context, email string) (*types.DeleteCustomerResult, error)
}
type service struct {
	userClient usersvc.Client
}

func NewService(userClient usersvc.Client) Service {
	return &service{
		userClient: userClient,
	}

}

func (s *service) CreateCustomer(ctx context.Context, payload types.CreateCustomerInput) (*userpb.CreateCustomerResponse, error) {
	req := &userpb.CreateCustomerRequest{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  payload.Password,
		Phone:     payload.Phone,
		Address:   payload.Address,
		AvatarUrl: payload.AvatarURL,
	}
	res, err := s.userClient.CreateCustomer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("user service create: %w", err)

	}
	return res, nil

}
