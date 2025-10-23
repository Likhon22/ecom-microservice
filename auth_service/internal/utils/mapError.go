package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	// tighten later to unwrap custom errors; placeholder maps to Internal
	return status.Errorf(codes.Internal, err.Error())
}
