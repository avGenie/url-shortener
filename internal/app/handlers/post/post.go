package post

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	db_err "github.com/avGenie/url-shortener/internal/app/storage/errors"
)

const (
	maxEncodedSize = 8
	timeout = 1*time.Second
)

type URLSaver interface {
	AddURL(ctx context.Context, key, value entity.URL) entity.Response
}

func postURLProcessing(saver URLSaver, ctx context.Context, inputURL, baseURIPrefix string) (string, error) {
	var shortURL *entity.URL

	userURL := entity.ParseURL(inputURL)
	added := false

	encodedURL := base64.StdEncoding.EncodeToString([]byte(inputURL))
	availableURLCount := len(encodedURL) / maxEncodedSize
	for i := 0; i < availableURLCount-1; i++ {
		shortURL = entity.ParseURL(encodedURL[(maxEncodedSize * i):(maxEncodedSize * (i + 1))])
		resp := saver.AddURL(ctx, *shortURL, *userURL)
		if resp.Status == entity.StatusOK {
			added = true
			break
		} else if !errors.Is(resp.Error, db_err.ErrURLAlreadyExists) {
			return "", resp.Error
		}
	}

	if !added {
		return "", nil
	}

	return fmt.Sprintf("%s/%s", baseURIPrefix, shortURL), nil
}