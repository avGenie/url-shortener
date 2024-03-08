package storage

import (
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/postgres"
)


func InitStorage(dbStorageConnect string) (entity.DBStorage, error) {
	var db entity.DBStorage
	db, err := postgres.New(dbStorageConnect)

	return db, err
}