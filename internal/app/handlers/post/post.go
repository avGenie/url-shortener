package post

import (
	"context"
	"encoding/hex"

	"errors"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/encoding"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"go.uber.org/zap"
)

const (
	maxEncodedSize = 8
	timeout        = 3 * time.Second
)

type URLSaver interface {
	SaveURL(ctx context.Context, key, value entity.URL) entity.URLResponse
}

type URLBatchSaver interface {
	SaveBatchURL(ctx context.Context, batch storage.Batch) storage.BatchResponse
}

func postURLProcessing(saver URLSaver, ctx context.Context, inputURL, baseURIPrefix string) (string, error) {
	hash := createHash(inputURL)
	if hash == "" {
		return "", fmt.Errorf("failed to create hash")
	}

	shortURL, err := entity.ParseURL(hash)
	if err != nil {
		zap.L().Error("error while parsing short url")
		return "", err
	}

	userURL, err := entity.ParseURL(inputURL)
	if err != nil {
		zap.L().Error("error while parsing user url")
		return "", err
	}

	resp := saver.SaveURL(ctx, *shortURL, *userURL)
	if resp.Status == entity.StatusError {
		if errors.Is(resp.Error, storage_err.ErrURLAlreadyExists) {
			return resp.URL.String(), resp.Error
		}
		return "", resp.Error
	}

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL.String()), nil
}

func createHash(url string) string {
	bs := encoding.NewSHA256([]byte(url))

	return hex.EncodeToString(bs)[:maxEncodedSize]
}

func createStorageBatch(urls models.ReqURLBatch) (storage.Batch, error) {
	dbBatch := make(storage.Batch, 0, len(urls))
	for _, url := range urls {
		shortURL := createHash(url.URL.String())
		if shortURL == "" {
			return nil, fmt.Errorf("failed to create hash")
		}

		obj := storage.BatchObject{
			ID:       url.Obj.ID,
			InputURL: url.URL.String(),
			ShortURL: shortURL,
		}

		dbBatch = append(dbBatch, obj)
	}

	return dbBatch, nil
}
