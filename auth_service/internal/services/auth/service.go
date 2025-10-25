package auth

import (
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/config"
	repo "auth_service/internal/repo/auth"
	"auth_service/internal/types"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateCustomer(ctx context.Context, payload types.CreateCustomerInput) (*userpb.CreateCustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*userpb.CreateCustomerResponse, error)
	// GetCustomers(ctx context.Context) ([]*types.CreateCustomerResult, error)
	// DeleteCustomer(ctx context.Context, email string) (*types.DeleteCustomerResult, error)
	Login(ctx context.Context, email, password, deviceID string) (string, string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) (string, error)
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

func (s *service) Login(ctx context.Context, email, password, deviceId string) (string, string, error) {
	credentials, err := s.userClient.GetCustomerCredentials(ctx, &userpb.GetCustomerByEmailRequest{Email: email})
	if err != nil {
		return "", "", fmt.Errorf("user service create: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
		return "", "", fmt.Errorf("password does not match: %w", err)
	}

	accessToken, err := utils.SignedToken(credentials.Email, credentials.Role, deviceId, s.authCnf.Jwt_Access_Token_Secret, s.authCnf.Access_Token_Exp_Duration)
	if err != nil {
		return "", "", fmt.Errorf("token error: %w", err)
	}

	refreshToken, err := utils.SignedToken(credentials.Email, credentials.Role, deviceId, s.authCnf.Jwt_Refresh_Token_Secret, s.authCnf.Refresh_Token_Exp_Duration)
	if err != nil {
		return "", "", fmt.Errorf("token error: %w", err)
	}

	if err := s.authRepo.Store(ctx, refreshToken, credentials.Email, deviceId, s.authCnf.Refresh_Token_Exp_Duration); err != nil {

		return "", "", fmt.Errorf("db error: %w", err)
	}
	log.Println("saved")

	return accessToken, refreshToken, nil
}

func (s *service) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {

	claims, err := utils.ParseJwt(refreshToken, s.authCnf.Jwt_Refresh_Token_Secret)
	if err != nil {
		return "", err

	}
	log.Println(claims)
	now := time.Now().UTC()
	tokenData, err := s.authRepo.Get(ctx, claims.Email, claims.DeviceId)
	if tokenData == nil {
		return "", errors.New("no token found")
	}

	if now.After(tokenData.ExpiresAt) {
		return "", errors.New("token expired")

	}

	if err != nil {
		return "", err
	}

	if tokenData.Token == "" {
		return "", errors.New("no token found")
	}

	if refreshToken != tokenData.Token {
		return "", errors.New("token doesnot match")
	}
	if tokenData.Revoked {
		return "", errors.New("token is revoked")

	}
	accessToken, err := utils.SignedToken(claims.Email, claims.Role, claims.DeviceId, s.authCnf.Jwt_Access_Token_Secret, s.authCnf.Access_Token_Exp_Duration)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.ParseJwt(refreshToken, s.authCnf.Jwt_Refresh_Token_Secret)
	if err != nil {
		return "", err
	}
	if err := s.authRepo.Revoked(ctx, claims.Email, claims.DeviceId); err != nil {
		return "", err
	}
	return "logout successful", nil
}
