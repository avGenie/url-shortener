package storage

import "sync"

// A thread-safe map which contains Add and Get methods
type Url struct {
	mutex sync.RWMutex
	urls  map[string]string
}

// Creates a new concurrent map
func NewUrl() *Url {
	var url Url
	url.urls = make(map[string]string)

	return &url
}

// Returns an element from the map
func (u *Url) Get(key string) (string, bool) {
	u.mutex.RLock()
	res, ok := u.urls[key]
	u.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
func (u *Url) Add(key, value string) {
	u.mutex.Lock()
	u.urls[key] = value
	u.mutex.Unlock()
}
