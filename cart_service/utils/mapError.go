package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		return status.Errorf(st.Code(), "%s", err.Error())
	}
	return status.Errorf(codes.Internal, "%s", err.Error())
}
