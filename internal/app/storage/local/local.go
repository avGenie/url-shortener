package local

import "github.com/avGenie/url-shortener/internal/app/entity"

type LocalStorage struct {
	urls map[entity.URL]entity.URL
}

func NewLocalStorage(size int) *LocalStorage {
	return &LocalStorage{
		urls: make(map[entity.URL]entity.URL, size),
	}
}

// Returns an element from the map
func (s *LocalStorage) Get(key entity.URL) (entity.URL, bool) {
	res, ok := s.urls[key]

	return res, ok
}

func (s *LocalStorage) GetAllURL() map[entity.URL]entity.URL {
	return s.urls
}

// Adds the given value under the specified key
func (s *LocalStorage) Add(key, value entity.URL) {
	s.urls[key] = value
}

// Adds the given value under the specified key
func (s *LocalStorage) Merge(inputStorage LocalStorage) {
	for key, value := range inputStorage.urls {
		s.urls[key] = value
	}
}
