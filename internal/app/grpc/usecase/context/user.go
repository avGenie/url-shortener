package context

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"google.golang.org/grpc/metadata"
)

const (
	userIDKey = "user_id"
)

func GetUserIDFromContext(ctx context.Context) entity.UserID {
	var userID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(userIDKey)
		if len(values) > 0 {
			userID = values[0]
		}
	}

	return entity.UserID(userID)
}

func SetUserIDContext(ctx context.Context, userID entity.UserID) context.Context {
	return metadata.AppendToOutgoingContext(ctx, userIDKey, userID.String())
}
