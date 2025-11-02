package authhandler

import (
	"auth_service/internal/services/auth"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
func (h *handler) CreateUserAccount(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.StandardResponse, error) {

	result, err := h.service.CreateUser(ctx, req)
	if err != nil {
		return nil, utils.MapError(err)

	}

	return utils.WrapSuccess(result, "user created successfully", 201), nil
}
func (h *handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.StandardResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, utils.MapError(err)
	}
	accessToken, refreshToken, err := h.service.Login(ctx, req.Email, req.Password, req.DeviceId)
	if err != nil {
		return nil, utils.MapError(err)

	}
	md := metadata.Pairs("set-cookie", utils.BuildCookieHeader("access-token", accessToken, 5*time.Minute, true, true))
	refreshmd := metadata.Pairs("set-cookie", utils.BuildCookieHeader("refresh-token", refreshToken, 24*time.Hour, true, true))
	grpc.SetHeader(ctx, md)
	grpc.SetHeader(ctx, refreshmd)
	resp := &userpb.LoginResponse{Message: refreshToken}
	return utils.WrapSuccess(resp, "login successful", 200), nil
}

func (h *handler) ValidateRefreshToken(ctx context.Context, req *userpb.ValidateRefreshTokenRequest) (*userpb.StandardResponse, error) {
	incomingMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapError(fmt.Errorf("no metadata in context"))
	}

	refreshToken := incomingMd.Get("refresh-token")

	accessToken, err := h.service.ValidateRefreshToken(ctx, refreshToken[0])
	if err != nil {
		return nil, utils.MapError(err)
	}
	md := metadata.Pairs("set-cookie", utils.BuildCookieHeader("access-token", accessToken, 5*time.Minute, true, true))
	grpc.SetHeader(ctx, md)
	resp := &userpb.ValidateRefreshTokenResponse{Message: accessToken}
	return utils.WrapSuccess(resp, "new access token generated", 200), nil

}

func (h *handler) Logout(ctx context.Context, req *userpb.LogoutRequest) (*userpb.StandardResponse, error) {
	incomingMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapError(fmt.Errorf("no metadata in context"))
	}
	refreshToken := incomingMetadata.Get("refresh-token")
	msg, err := h.service.Logout(ctx, refreshToken[0])
	if err != nil {
		return nil, utils.MapError(err)

	}
	resp := &userpb.LogoutResponse{Message: msg}
	return utils.WrapSuccess(resp, "logout successful", 200), nil

}
