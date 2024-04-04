package model

import (
	"context"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
)

type Storage interface {
	Close() entity.Response
	PingServer(ctx context.Context) error

	SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error
	SaveBatchURL(ctx context.Context, userID entity.UserID, batch Batch) (Batch, error)

	GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error)
	GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error)

	DeleteBatchURL(ctx context.Context, urls entity.DeletedURLBatch) error
}
