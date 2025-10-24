package utils

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

func GetAccessTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	// gRPC metadata keys are lowercase by convention
	authHeaders := md.Get("set-cookie") // your token is in "set-cookie" key

	for _, header := range authHeaders {
		// Example: access-token=xxx; Path=/; Max-Age=300; HttpOnly; Secure; SameSite=Strict
		if strings.HasPrefix(header, "refresh-token=") {
			parts := strings.Split(header, ";")
			tokenPair := strings.SplitN(parts[0], "=", 2)
			if len(tokenPair) == 2 {
				return tokenPair[1], nil
			}
		}
	}

	return "", fmt.Errorf("access token not found")
}
