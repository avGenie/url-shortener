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
func (s *TSLocalStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	s.mutex.RLock()
	res, ok := s.urls.Get(userID, key)
	s.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error while getting url from ts local storage: %w", api.ErrShortURLNotFound)
	}

	return &res, nil
}

func (s *TSLocalStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	s.mutex.RLock()
	urls, ok := s.urls.GetByUserID(userID)
	s.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error while getting all user url from ts local storage: %w", api.ErrShortURLNotFound)
	}

	var allURLs models.AllUrlsBatch
	for key, value := range urls {
		allURLs = append(allURLs, models.AllUrlsResponse{
			ShortURL: key.String(),
			OriginalURL: value.String(),
		})
	}

	return allURLs, nil
}

// Adds the given value under the specified key
func (s *TSLocalStorage) SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.urls.Get(userID, key)
	if ok {
		return fmt.Errorf("error while save url to ts local storage: %w", api.ErrURLAlreadyExists)
	}

	s.urls.Add(userID, key, value)

	return nil
}

// Adds elements from the given batch to the local storage
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

		localUrls.Add(userID, *key, *value)
	}

	s.mutex.Lock()
	s.urls.Merge(*localUrls)
	s.mutex.Unlock()

	return batch, nil
}

func (s *TSLocalStorage) AddUser(ctx context.Context, userID entity.UserID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.urls.GetByUserID(userID)
	if ok {
		return fmt.Errorf("error while save url to ts local storage: %w", api.ErrUserAlreadyExists)
	}

	s.urls.AddUserID(userID)

	return nil
}

func (s *TSLocalStorage) AuthUser(ctx context.Context, userID entity.UserID) (entity.UserID, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.urls.GetByUserID(userID)
	if !ok {
		return "", fmt.Errorf("error while save url to ts local storage: %w", api.ErrUserIDNotFound)
	}

	return userID, nil
}

func PingServer(ctx context.Context) error {
	return nil
}

func Close() entity.Response {
	return entity.OKResponse()
}
