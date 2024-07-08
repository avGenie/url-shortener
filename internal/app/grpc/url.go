package grpc

import (
	"context"
	"fmt"

	grpc_context "github.com/avGenie/url-shortener/internal/app/grpc/usecase/context"
	pb "github.com/avGenie/url-shortener/proto"
)

func (s *ShortenerServer) GetShortURL(ctx context.Context, original *pb.OriginalURL) (*pb.ShortURL, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)
	fmt.Println(userID)

	return &pb.ShortURL{}, nil
}

func (s *ShortenerServer) GetOriginalURL(ctx context.Context, original *pb.ShortURL) (*pb.OriginalURL, error) {
	userID := grpc_context.GetUserIDFromContext(ctx)
	fmt.Println(userID)

	return &pb.OriginalURL{}, nil
}
