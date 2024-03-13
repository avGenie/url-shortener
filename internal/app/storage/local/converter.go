package local

import (
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
)

func CreateLocalStorageFromBatch(batch model.Batch) (*LocalStorage, error) {
	localUrls := NewLocalStorage(len(batch))
	for _, obj := range batch {
		key, err := entity.NewURL(obj.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create input url from batch in local storage: %w", err)
		}
		value, err := entity.NewURL(obj.InputURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create short url from batch in local storage: %w", err)
		}

		localUrls.Add(*key, *value)
	}

	return localUrls, nil
}