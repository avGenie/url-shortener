package storage

import (
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/storage/file"
	"github.com/avGenie/url-shortener/internal/app/storage/local"
	"github.com/avGenie/url-shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
)

func InitStorage(config config.Config) (entity.Storage, error) {
	var db entity.Storage
	var err error

	if len(config.DBStorageConnect) > 0 {
		zap.L().Info("init postgres storage")
		db, err = postgres.NewPostgresStorage(config.DBStorageConnect)
	} else if len(config.DBFileStoragePath) > 0 {
		zap.L().Info("init file storage")
		db, err = file.NewFileStorage(config.DBFileStoragePath)
	} else {
		zap.L().Info("init local storage")
		db = local.NewTsLocalStorage(0)
		err = nil
	}

	return db, err
}
