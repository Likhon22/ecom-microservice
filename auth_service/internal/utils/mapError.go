package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {

	return status.Errorf(codes.Internal, err.Error())
}
