package model

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

type Storage interface {
	Close() entity.Response
	PingServer(ctx context.Context) error
	SaveURL(ctx context.Context, key, value entity.URL) error
	SaveBatchURL(ctx context.Context, batch Batch) (Batch, error)
	GetURL(ctx context.Context, key entity.URL) (*entity.URL, error)
}
