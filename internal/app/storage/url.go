package storage

import (
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

// A thread-safe map which contains Add and Get methods
type URLStorage struct {
	mutex sync.RWMutex
	urls  map[entity.URL]entity.URL
}

// Creates a new concurrent map
func NewURLStorage() *URLStorage {
	var url URLStorage
	url.urls = make(map[entity.URL]entity.URL)

	return &url
}

// Returns an element from the map
func (u *URLStorage) Get(key entity.URL) (entity.URL, bool) {
	u.mutex.RLock()
	res, ok := u.urls[key]
	u.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
func (u *URLStorage) Add(key, value entity.URL) {
	u.mutex.Lock()
	u.urls[key] = value
	u.mutex.Unlock()
}
