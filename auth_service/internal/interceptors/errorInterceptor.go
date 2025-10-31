package interceptors

import (
	userpb "auth_service/proto/gen"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterCeptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		resp, err := handler(ctx, req)
		if err != nil {
			statusCode := grpcCodeToHTTP(status.Code(err))
			errorResponse := &userpb.StandardResponse{
				Success:    false,
				Message:    err.Error(),
				StatusCode: int32(statusCode),
				Result:     nil,
			}
			return errorResponse, nil
		}
		return resp, nil
	}

}
func grpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return 400
	case codes.Unauthenticated:
		return 401
	case codes.PermissionDenied:
		return 403
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.Internal:
		return 500
	default:
		return 500
	}
}
