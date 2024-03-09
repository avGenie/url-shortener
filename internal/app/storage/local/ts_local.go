package local

import (
	"context"
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/errors"
)

type TsLocalStorage struct {
	entity.Storage

	mutex sync.RWMutex
	urls LocalStorage
}

func NewTsLocalStorage(size int) *TsLocalStorage {
	return &TsLocalStorage{
		urls: *NewLocalStorage(size),
	}
}

// Returns an element from the map
func (s *TsLocalStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return entity.ErrorURLResponse(errors.ErrShortURLNotFound)
	}

	return entity.OKURLResponse(res)
}

// Adds the given value under the specified key
func (s *TsLocalStorage) AddURL(ctx context.Context, key, value entity.URL) entity.Response {
	s.mutex.Lock()
	s.urls.Add(key, value)
	s.mutex.Unlock()

	return entity.OKResponse()
}

func PingServer(ctx context.Context) entity.Response {
	return entity.OKResponse()
}

func Close() entity.Response {
	return entity.OKResponse()
}