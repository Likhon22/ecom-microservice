package auth

import (
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/types"
	userpb "auth_service/proto/gen"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateCustomer(ctx context.Context, payload types.CreateCustomerInput) (*userpb.CreateCustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*userpb.CreateCustomerResponse, error)
	// GetCustomers(ctx context.Context) ([]*types.CreateCustomerResult, error)
	// DeleteCustomer(ctx context.Context, email string) (*types.DeleteCustomerResult, error)
	Login(ctx context.Context, email, password string) (string, error)
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
func (s *service) GetCustomerByEmail(ctx context.Context, email string) (*userpb.CreateCustomerResponse, error) {
	req := &userpb.GetCustomerByEmailRequest{
		Email: email,
	}
	res, err := s.userClient.GetCustomerByEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("user service create: %w", err)

	}
	return res, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	credentials, err := s.userClient.GetCustomerCredentials(ctx, &userpb.GetCustomerByEmailRequest{Email: email})
	if err != nil {
		return "", fmt.Errorf("user service create: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {

		return "", fmt.Errorf("password does not match: %w", err)
	}
	return "successfully login", nil
}
