package post

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"errors"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"go.uber.org/zap"
)

const (
	maxEncodedSize = 8
	timeout        = 3 * time.Second
)

type URLSaver interface {
	AddURL(ctx context.Context, key, value entity.URL) entity.Response
}

func postURLProcessing(saver URLSaver, ctx context.Context, inputURL, baseURIPrefix string) (string, error) {
	h := sha256.New()
	h.Write([]byte(inputURL))
	bs := h.Sum(nil)

	shortURL, err := entity.ParseURL(hex.EncodeToString(bs)[:maxEncodedSize])
	if err != nil {
		zap.L().Error("error while parsing short url")
		return "", err
	}

	userURL, err := entity.ParseURL(inputURL)
	if err != nil {
		zap.L().Error("error while parsing user url")
		return "", err
	}

	resp := saver.AddURL(ctx, *shortURL, *userURL)
	if resp.Status == entity.StatusError {
		if !errors.Is(resp.Error, storage.ErrURLAlreadyExists) {
			return "", resp.Error
		}
		return "", resp.Error
	}

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL), nil
}
