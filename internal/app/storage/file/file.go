// Package file contains implementation of file storage
package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/avGenie/url-shortener/internal/app/storage/local"
)

// FileStorage File storage object
type FileStorage struct {
	model.Storage

	mutex sync.RWMutex

	file     *os.File
	fileName string
	encoder  *json.Encoder

	cache  local.LocalStorage
	lastID uint

	IsTemp bool
}

// NewFileStorage Creates a new file storage object
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
		mutex:    sync.RWMutex{},
		file:     file,
		fileName: fileName,
		encoder:  json.NewEncoder(file),
		cache:    *local.NewLocalStorage(0),
		lastID:   0,
	}

	err = storage.fillCacheFromFile()
	if err != nil {
		return nil, err
	}

	zap.L().Info("storage was created successfully")

	return storage, nil
}

// GetURL Returns user URL from file storage
func (s *FileStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	s.mutex.RLock()
	if s.file == nil {
		return nil, fmt.Errorf("error while getting url from file: %w", api.ErrFileStorageNotOpen)
	}
	res, ok := s.cache.Get(key)
	s.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("error while getting url from file: %w", api.ErrShortURLNotFound)
	}

	return &res, nil
}

// GetAllURLByUserID Returns all user URL from file storage
func (s *FileStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	s.mutex.RLock()
	if s.file == nil {
		return nil, fmt.Errorf("error while getting url from file: %w", api.ErrFileStorageNotOpen)
	}
	urls := s.cache.GetAllURL()
	s.mutex.RUnlock()

	var allURLs models.AllUrlsBatch
	for key, value := range urls {
		allURLs = append(allURLs, models.AllUrlsResponse{
			ShortURL:    key.String(),
			OriginalURL: value.String(),
		})
	}

	return allURLs, nil
}

// SaveURL Saves user URL to file storage
func (s *FileStorage) SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.file == nil {
		return fmt.Errorf("error while save url to file storage: %w", api.ErrFileStorageNotOpen)
	}

	storageRec := &entity.URLRecord{
		ID:          s.lastID + 1,
		ShortURL:    key.Path,
		OriginalURL: value.String(),
	}

	err := s.encoder.Encode(&storageRec)
	if err != nil {
		return fmt.Errorf("error while encoding entity for file commit: %w", err)
	}

	s.file.Sync()

	s.cache.Add(key, value)
	s.lastID = storageRec.ID

	return nil
}

// SaveBatchURL Saves batch of user URLs to file storage
func (s *FileStorage) SaveBatchURL(ctx context.Context, userID entity.UserID, batch model.Batch) (model.Batch, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.file == nil {
		return nil, api.ErrFileStorageNotOpen
	}

	localUrls := local.NewLocalStorage(len(batch))
	records := make([]entity.URLRecord, 0, len(batch))
	for _, obj := range batch {
		key, err := entity.NewURL(obj.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("exit to create input url from batch in file storage: %w", err)
		}
		value, err := entity.NewURL(obj.InputURL)
		if err != nil {
			return nil, fmt.Errorf("exit to create short url from batch in file storage: %w", err)
		}

		localUrls.Add(*key, *value)

		storageRec := &entity.URLRecord{
			ID:          s.lastID + 1,
			ShortURL:    key.Path,
			OriginalURL: value.String(),
		}

		s.lastID = storageRec.ID

		records = append(records, *storageRec)
	}

	err := s.encoder.Encode(&records)
	if err != nil {
		return nil, fmt.Errorf("error while encoding entity for file commit: %w", err)
	}
	s.file.Sync()

	s.cache.Merge(*localUrls)

	return batch, nil
}

// Close Closes connection to file storage
func (s *FileStorage) Close() {
	s.file.Name()
	if strings.Contains(s.fileName, os.TempDir()) {
		err := os.Remove(s.fileName)
		if err != nil {
			zap.L().Error("error while closing file storage", zap.Error(err))
		}
	}
}

// PingServer Pings to file storage
func (s *FileStorage) PingServer(ctx context.Context) error {
	if s.file == nil {
		return fmt.Errorf("error while ping file storage: %w", api.ErrFileStorageNotOpen)
	}

	return nil
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
