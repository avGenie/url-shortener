package storage

import (
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/postgres"
)

func InitStorage(dbStorageConnect string) (entity.Storage, error) {
	var db entity.Storage
	db, err := postgres.New(dbStorageConnect)

	return db, err
}
