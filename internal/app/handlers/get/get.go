package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"go.uber.org/zap"
)

const (
	timeout     = 3 * time.Second
	pingTimeout = 1 * time.Second
)

// Errors returning while GET request processing
//
// ErrInternal - internal error
// ErrAllURLNotFound - returned if all URL is not found in storage for given user
var (
	ErrInternal       = errors.New("getting internal error while getting all user urls")
	ErrAllURLNotFound = errors.New("urls for this user not found")
)

// ProcessAllUserURL Returns all URLs for given user
func ProcessAllUserURL(getter AllURLGetter, ctx context.Context, userID entity.UserID, baseURIPrefix string) ([]byte, error) {
	urls, err := getter.GetAllURLByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("couldn't get all user urls", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, ErrInternal
	}

	if len(urls) == 0 {
		return nil, ErrAllURLNotFound
	}

	for index, url := range urls {
		url.ShortURL = fmt.Sprintf("%s/%s", baseURIPrefix, url.ShortURL)
		urls[index] = url
	}

	out, err := json.Marshal(urls)
	if err != nil {
		zap.L().Error("error while converting all user urls to output", zap.Error(err))
		return nil, ErrInternal
	}

	return out, nil
}
