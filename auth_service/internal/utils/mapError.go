package utils

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	st, ok := status.FromError(err)
	log.Println("code", st.Code())
	log.Println("error", st.Message())
	if ok {
		return status.Errorf(st.Code(), "%s", st.Message())
	}
	return status.Errorf(codes.Internal, "%s", err.Error())
}
