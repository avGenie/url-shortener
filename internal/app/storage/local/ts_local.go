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
func (s *TSLocalStorage) GetURL(ctx context.Context, key entity.URL) (*entity.URL, error) {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error while getting url from ts local storage: %w", api.ErrShortURLNotFound)
	}

	return &res, nil
}

// Adds the given value under the specified key
func (s *TSLocalStorage) SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.urls.Get(key)
	if ok {
		return fmt.Errorf("error while save url to ts local storage: %w", api.ErrURLAlreadyExists)
	}

	s.urls.Add(key, value)

	return nil
}

// Adds elements from the given batch to the local storage
func (s *TSLocalStorage) SaveBatchURL(ctx context.Context, batch model.Batch) (model.Batch, error) {
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

	s.mutex.Lock()
	s.urls.Merge(*localUrls)
	s.mutex.Unlock()

	return batch, nil
}

func PingServer(ctx context.Context) error {
	return nil
}

func Close() entity.Response {
	return entity.OKResponse()
}
