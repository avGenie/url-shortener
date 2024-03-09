package file

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/local"
	"go.uber.org/zap"
)

type FileStorage struct {
	mutex sync.RWMutex

	file    *os.File
	encoder *json.Encoder

	cache  local.LocalStorage
	lastID uint

	IsTemp bool
}

// Creates a new concurrent map
func NewFileStorage(fileName string) (*FileStorage, error) {
	if fileName == "" {
		zap.L().Info("storage was created successfully without keeping URL on disk")
		return &FileStorage{
			cache:  *local.NewLocalStorage(0),
			lastID: 0,
		}, nil
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := &FileStorage{
		mutex:   sync.RWMutex{},
		file:    file,
		encoder: json.NewEncoder(file),
		cache:   *local.NewLocalStorage(0),
		lastID:  0,
	}

	err = storage.fillCacheFromFile()
	if err != nil {
		return nil, err
	}

	zap.L().Info("storage was created successfully")

	return storage, nil
}

// Returns an element from the map
func (s *FileStorage) Get(key entity.URL) (entity.URL, bool) {
	s.mutex.RLock()
	res, ok := s.cache.Get(key)
	s.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
//
// Returns `true` if element has been added to the storage.
func (s *FileStorage) Add(key, value entity.URL) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	
	if _, ok := s.cache.Get(key); ok {
		return false, nil
	}

	if s.file == nil {
		s.cache.Add(key, value)
		s.lastID++
		return true, nil
	}

	storageRec := &entity.URLRecord{
		ID:          s.lastID + 1,
		ShortURL:    key.Path,
		OriginalURL: value.String(),
	}

	err := s.encoder.Encode(&storageRec)
	if err != nil {
		return false, err
	}

	s.file.Sync()

	s.cache.Add(key, value)
	s.lastID = storageRec.ID

	return true, nil
}

// Fills cache from the DB storage file
func (s *FileStorage) fillCacheFromFile() error {
	s.file.Seek(0, 0)
	scanner := bufio.NewScanner(s.file)
	record := entity.URLRecord{}

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return err
		}

		key, err := entity.NewURL(record.ShortURL)
		if err != nil {
			return err
		}

		value, err := entity.NewURL(record.OriginalURL)
		if err != nil {
			return err
		}

		s.cache.Add(*key, *value)
		s.lastID = record.ID
	}

	return nil
}
