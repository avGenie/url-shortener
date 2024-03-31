package local

import "github.com/avGenie/url-shortener/internal/app/entity"

type LocalStorage struct {
	urls map[entity.UserID]map[entity.URL]entity.URL
}

func NewLocalStorage(size int) *LocalStorage {
	return &LocalStorage{
		urls: make(map[entity.UserID]map[entity.URL]entity.URL, size),
	}
}

// Returns an element from the map
func (s *LocalStorage) Get(userID entity.UserID, key entity.URL) (entity.URL, bool) {
	res, ok := s.urls[userID][key]

	return res, ok
}

// Adds the given value under the specified key
func (s *LocalStorage) Add(userID entity.UserID, key, value entity.URL) {
	_, ok := s.GetByUserID(userID)
	if !ok {
		s.AddUserID(userID)
	}
	
	s.urls[userID][key] = value
}

// Adds user id to the storage
func (s *LocalStorage) AddUserID(userID entity.UserID) {
	s.urls[userID] = make(map[entity.URL]entity.URL)
}

// Adds user id to the storage
func (s *LocalStorage) GetByUserID(userID entity.UserID) (map[entity.URL]entity.URL, bool) {
	res, ok := s.urls[userID]

	return res, ok
}

// Adds the given value under the specified key
func (s *LocalStorage) Merge(inputStorage LocalStorage) {
	for key, value := range inputStorage.urls {
		s.urls[key] = value
	}
}
