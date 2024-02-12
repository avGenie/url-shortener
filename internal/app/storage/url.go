package storage

import "sync"

// A thread-safe map which contains Add and Get methods
type URLStorage struct {
	mutex sync.RWMutex
	urls  map[string]string
}

// Creates a new concurrent map
func NewURLStorage() *URLStorage {
	var url URLStorage
	url.urls = make(map[string]string)

	return &url
}

// Returns an element from the map
func (u *URLStorage) Get(key string) (string, bool) {
	u.mutex.RLock()
	res, ok := u.urls[key]
	u.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
func (u *URLStorage) Add(key, value string) {
	u.mutex.Lock()
	u.urls[key] = value
	u.mutex.Unlock()
}
