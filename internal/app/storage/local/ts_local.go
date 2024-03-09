package local

import (
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

type TsLocalStorage struct {
	mutex sync.RWMutex
	urls LocalStorage
}

func NewTsLocalStorage(size int) *TsLocalStorage {
	return &TsLocalStorage{
		urls: *NewLocalStorage(size),
	}
}

// Returns an element from the map
func (s *TsLocalStorage) Get(key entity.URL) (entity.URL, bool) {
	s.mutex.RLock()
	res, ok := s.urls.Get(key)
	s.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
func (s *TsLocalStorage) Add(key, value entity.URL) {
	s.mutex.Lock()
	s.urls.Add(key, value)
	s.mutex.Unlock()
}