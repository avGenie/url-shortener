package model

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

type Storage interface {
	Close() entity.Response
	PingServer(ctx context.Context) entity.Response
	SaveURL(ctx context.Context, key, value entity.URL) entity.Response
	SaveBatchURL(ctx context.Context, batch Batch) BatchResponse
	GetURL(ctx context.Context, key entity.URL) entity.URLResponse
}
