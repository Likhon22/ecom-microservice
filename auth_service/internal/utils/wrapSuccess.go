package utils

import userpb "auth_service/proto/gen"

func WrapSuccess(result interface{}, message string, statusCode int32) *userpb.StandardResponse {
	switch v := result.(type) {
	case *userpb.LoginResponse:
		return &userpb.StandardResponse{
			Success:    true,
			Message:    message,
			StatusCode: statusCode,
			Result:     &userpb.StandardResponse_LoginData{LoginData: v},
		}
	case *userpb.LogoutResponse:
		return &userpb.StandardResponse{
			Success:    true,
			Message:    message,
			StatusCode: statusCode,
			Result:     &userpb.StandardResponse_LogoutData{LogoutData: v},
		}
	case *userpb.ValidateRefreshTokenResponse:
		return &userpb.StandardResponse{
			Success:    true,
			Message:    message,
			StatusCode: statusCode,
			Result:     &userpb.StandardResponse_RefreshTOkenData{RefreshTOkenData: v},
		}
	case *userpb.CreateUserResponse:
		return &userpb.StandardResponse{
			Success:    true,
			Message:    message,
			StatusCode: statusCode,
			Result:     &userpb.StandardResponse_UserData{UserData: v},
		}
	default:
		return &userpb.StandardResponse{
			Success:    true,
			Message:    message,
			StatusCode: statusCode,
		}
	}
}
