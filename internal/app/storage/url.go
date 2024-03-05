package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/logger"
)

type URLStorage struct {
	mutex sync.RWMutex

	file    *os.File
	encoder *json.Encoder

	cache  map[entity.URL]entity.URL
	lastID uint

	IsTemp bool
}

// Creates a new concurrent map
func NewURLStorage(fileName string) (*URLStorage, error) {
	if fileName == "" {
		logger.Log.Info("storage was created successfully without keeping URL on disk")
		return &URLStorage{
			cache:  make(map[entity.URL]entity.URL),
			lastID: 0,
		}, nil
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := &URLStorage{
		mutex:   sync.RWMutex{},
		file:    file,
		encoder: json.NewEncoder(file),
		cache:   make(map[entity.URL]entity.URL),
		lastID:  0,
	}

	err = storage.fillCacheFromFile()
	if err != nil {
		return nil, err
	}

	logger.Log.Info("storage was created successfully")

	return storage, nil
}

// Fills cache from the DB storage file
func (u *URLStorage) fillCacheFromFile() error {
	u.file.Seek(0, 0)
	scanner := bufio.NewScanner(u.file)
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

		u.cache[*key] = *value
		u.lastID = record.ID
	}

	return nil
}

// Returns an element from the map
func (u *URLStorage) Get(key entity.URL) (entity.URL, bool) {
	u.mutex.RLock()
	res, ok := u.cache[key]
	u.mutex.RUnlock()

	return res, ok
}

// Adds the given value under the specified key
func (u *URLStorage) Add(key, value entity.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.file == nil {
		u.cache[key] = value
		u.lastID++
		return nil
	}

	storageRec := &entity.URLRecord{
		ID:          u.lastID + 1,
		ShortURL:    key.Path,
		OriginalURL: value.String(),
	}

	err := u.encoder.Encode(&storageRec)
	if err != nil {
		return err
	}

	u.file.Sync()

	u.cache[key] = value
	u.lastID = storageRec.ID

	return nil
}
