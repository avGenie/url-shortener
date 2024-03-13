package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
)

type TSLocalStorage struct {
	model.Storage

	mutex sync.RWMutex
	urls  LocalStorage
}

func NewTSLocalStorage(size int) *TSLocalStorage {
	return &TSLocalStorage{
		urls: *NewLocalStorage(size),
	}
}

// Returns an element from the map
func (s *TSLocalStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return entity.ErrorURLResponse(api.ErrShortURLNotFound)
	}

	return entity.OKURLResponse(res)
}

// Adds the given value under the specified key
func (s *TSLocalStorage) SaveURL(ctx context.Context, key, value entity.URL) entity.URLResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, ok := s.urls.Get(key)
	if ok {
		return entity.ErrorURLValueResponse(api.ErrURLAlreadyExists, res)
	}

	s.urls.Add(key, value)

	return entity.OKURLResponse(entity.URL{})
}

// Adds elements from the given batch to the local storage
func (s *TSLocalStorage) SaveBatchURL(ctx context.Context, batch model.Batch) model.BatchResponse {
	localUrls := NewLocalStorage(len(batch))
	for _, obj := range batch {
		key, err := entity.NewURL(obj.ShortURL)
		if err != nil {
			return model.ErrorBatchResponse(fmt.Errorf("failed to create input url from batch in local storage: %w", err))
		}
		value, err := entity.NewURL(obj.InputURL)
		if err != nil {
			return model.ErrorBatchResponse(fmt.Errorf("failed to create short url from batch in local storage: %w", err))
		}

		localUrls.Add(*key, *value)
	}

	s.mutex.Lock()
	s.urls.Merge(*localUrls)
	s.mutex.Unlock()

	return model.OKBatchResponse(batch)
}

func PingServer(ctx context.Context) entity.Response {
	return entity.OKResponse()
}

func Close() entity.Response {
	return entity.OKResponse()
}
