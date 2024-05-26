package handlers

import (
	"context"
	"encoding/hex"

	"errors"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/converter"
	"github.com/avGenie/url-shortener/internal/app/encoding"
	"github.com/avGenie/url-shortener/internal/app/entity"
	post_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
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
	SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error
}

type URLBatchSaver interface {
	SaveBatchURL(ctx context.Context, userID entity.UserID, batch storage.Batch) (storage.Batch, error)
}

func postURLProcessing(saver URLSaver, ctx context.Context, userID entity.UserID,
	inputURL, baseURIPrefix string) (string, error) {
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

	err = saver.SaveURL(ctx, userID, *shortURL, *userURL)
	if err != nil {
		if errors.Is(err, storage_err.ErrURLAlreadyExists) {
			return createOutputPostString(baseURIPrefix, shortURL.String()), err
		}
		return "", err
	}

	return createOutputPostString(baseURIPrefix, shortURL.String()), nil
}

func createOutputPostString(baseURIPrefix, url string) string {
	return fmt.Sprintf("%s/%s", baseURIPrefix, url)
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

func batchURLProcessing(saver URLBatchSaver, ctx context.Context, userID entity.UserID,
	batch models.ReqBatch, baseURIPrefix string) (models.ResBatch, error) {
	urls, err := converter.ConvertBatchReqToURL(batch)
	if err != nil {
		zap.L().Error(post_err.CannotProcessURL, zap.Error(err))
		return nil, fmt.Errorf(post_err.WrongJSONFormat)
	}

	sBatch, err := createStorageBatch(urls)
	if err != nil {
		zap.L().Error("error while creating storage batch", zap.Error(err))
		return nil, fmt.Errorf(post_err.InternalServerError)
	}

	savedBatch, err := saver.SaveBatchURL(ctx, userID, sBatch)
	if err != nil {
		zap.L().Error("error while saving url to storage", zap.Error(err))
		return nil, fmt.Errorf(post_err.InternalServerError)
	}

	outBatch := converter.ConvertStorageBatchToOutBatch(savedBatch, baseURIPrefix)

	return outBatch, nil
}
