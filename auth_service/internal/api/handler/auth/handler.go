package authhandler

import (
	"auth_service/internal/services/auth"
	"auth_service/internal/utils"
	userpb "auth_service/proto/gen"
	"context"
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
func (h *handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	accessToken, refreshToken, err := h.service.Login(ctx, req.Email, req.Password, req.DeviceId)
	if err != nil {
		return nil, utils.MapError(err)

	}
	md := metadata.Pairs("set-cookie", utils.BuildCookieHeader("access-token", accessToken, 5*time.Minute, true, true))
	refreshmd := metadata.Pairs("set-cookie", utils.BuildCookieHeader("refresh-token", refreshToken, 24*time.Hour, true, true))
	grpc.SetHeader(ctx, md)
	grpc.SetHeader(ctx, refreshmd)
	return &userpb.LoginResponse{Message: accessToken}, nil
}

func (h *handler) ValidateRefreshToken(ctx context.Context, req *userpb.ValidateRefreshTokenRequest) (*userpb.ValidateRefreshTokenResponse, error) {
	token, err := utils.GetAccessTokenFromCtx(ctx)
	if err != nil {
		return nil, utils.MapError(err)
	}
}
