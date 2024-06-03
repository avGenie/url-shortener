package local

import "github.com/avGenie/url-shortener/internal/app/entity"

// LocalStorage Local storage object
type LocalStorage struct {
	urls map[entity.URL]entity.URL
}

// NewLocalStorage Creates local storage object
func NewLocalStorage(size int) *LocalStorage {
	return &LocalStorage{
		urls: make(map[entity.URL]entity.URL, size),
	}
}

// Get Returns an element from the local storage
func (s *LocalStorage) Get(key entity.URL) (entity.URL, bool) {
	res, ok := s.urls[key]

	return res, ok
}

// GetAllURL Returns all storage
func (s *LocalStorage) GetAllURL() map[entity.URL]entity.URL {
	return s.urls
}

// Add Adds the given value under the specified key to local storage
func (s *LocalStorage) Add(key, value entity.URL) {
	s.urls[key] = value
}

// Merge Adds the given value under the specified key to local storage
func (s *LocalStorage) Merge(inputStorage LocalStorage) {
	for key, value := range inputStorage.urls {
		s.urls[key] = value
	}
}
