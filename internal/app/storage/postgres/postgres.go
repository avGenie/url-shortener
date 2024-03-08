package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/storage"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	storage.DBStorage

	db *sql.DB
}

func New(dbStorageConnect string) (*PostgresStorage, error) {
	db ,err := sql.Open("postgres", dbStorageConnect)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql connect: %w", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s *PostgresStorage) PingDBServer(ctx context.Context) int {
	err := s.db.PingContext(ctx)
	if err != nil {
		zap.L().Error("cannot ping postgres db", zap.String("error", err.Error()))
        return http.StatusInternalServerError
    }

	return http.StatusOK
}