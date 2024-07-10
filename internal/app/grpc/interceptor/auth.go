package interceptor

import (
	"context"

	grpc_context "github.com/avGenie/url-shortener/internal/app/grpc/usecase/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthInterceptor Checks user id from context
func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)

	if len(userID) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	return handler(ctx, req)
}
