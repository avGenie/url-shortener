// Package local contains implementation of local storage
package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
)

// TSLocalStorage Thread save local storage object
type TSLocalStorage struct {
	model.Storage

	mutex sync.RWMutex
	urls  LocalStorage
}

// NewTSLocalStorage Creates thread save local storage object
func NewTSLocalStorage(size int) *TSLocalStorage {
	return &TSLocalStorage{
		urls: *NewLocalStorage(size),
	}
}

// GetURL Returns URL user from local storage
func (s *TSLocalStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error while getting url from ts local storage: %w", api.ErrShortURLNotFound)
	}

	return &res, nil
}

// GetAllURLByUserID Returns all user URLs from local storage
func (s *TSLocalStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	s.mutex.RLock()
	urls := s.urls.GetAllURL()
	s.mutex.RUnlock()

	var allURLs models.AllUrlsBatch
	for key, value := range urls {
		allURLs = append(allURLs, models.AllUrlsResponse{
			ShortURL:    key.String(),
			OriginalURL: value.String(),
		})
	}

	return allURLs, nil
}

// SaveURL Saves user URL to local storage
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

// SaveBatchURL Saves batch of user URLs to local storage
func (s *TSLocalStorage) SaveBatchURL(ctx context.Context, userID entity.UserID, batch model.Batch) (model.Batch, error) {
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

// PingServer Pings to local storage
func (s *TSLocalStorage) PingServer(ctx context.Context) error {
	return nil
}

// Close Closes connection to local storage
func (s *TSLocalStorage) Close() {
}
