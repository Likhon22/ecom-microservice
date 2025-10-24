package auth

import (
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/config"
	repo "auth_service/internal/repo/auth"
	"auth_service/internal/types"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateCustomer(ctx context.Context, payload types.CreateCustomerInput) (*userpb.CreateCustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*userpb.CreateCustomerResponse, error)
	// GetCustomers(ctx context.Context) ([]*types.CreateCustomerResult, error)
	// DeleteCustomer(ctx context.Context, email string) (*types.DeleteCustomerResult, error)
	Login(ctx context.Context, email, password, deviceID string) (string, error)
}
type service struct {
	userClient usersvc.Client
	authCnf    *config.AuthConfig
	authRepo   repo.AuthRepo
}

func NewService(userClient usersvc.Client, authCnf *config.AuthConfig, authRepo repo.AuthRepo) Service {
	return &service{
		userClient: userClient,
		authCnf:    authCnf,
		authRepo:   authRepo,
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

func (s *service) Login(ctx context.Context, email, password, deviceId string) (string, error) {
	credentials, err := s.userClient.GetCustomerCredentials(ctx, &userpb.GetCustomerByEmailRequest{Email: email})
	if err != nil {
		return "", fmt.Errorf("user service create: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {

		return "", fmt.Errorf("password does not match: %w", err)
	}

	accessToken, err := utils.SignedToken(credentials.Email, credentials.Role, s.authCnf.Jwt_Access_Token_Secret, s.authCnf.Access_Token_Exp_Duration)
	if err != nil {
		return "", fmt.Errorf("token error: %w", err)
	}

	refreshToken, err := utils.SignedToken(credentials.Email, credentials.Role, s.authCnf.Jwt_Refresh_Token_Secret, s.authCnf.Refresh_Token_Exp_Duration)
	if err != nil {
		return "", fmt.Errorf("token error: %w", err)
	}

	if err := s.authRepo.Store(ctx, refreshToken, credentials.Email, deviceId, s.authCnf.Refresh_Token_Exp_Duration); err != nil {

		return "", fmt.Errorf("db error: %w", err)
	}
	log.Println("saved")
	wholeToken := fmt.Sprintf("refreshToken: %s accessToken: %s", refreshToken, accessToken)
	return wholeToken, nil
}
