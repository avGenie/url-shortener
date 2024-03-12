package model

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

type Storage interface {
	Close() entity.Response
	PingServer(ctx context.Context) entity.Response
	SaveURL(ctx context.Context, key, value entity.URL) entity.Response
	GetURL(ctx context.Context, key entity.URL) entity.URLResponse
}
