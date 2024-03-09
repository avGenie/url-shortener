package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	entity.Storage

	db *sql.DB
}

func New(dbStorageConnect string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dbStorageConnect)
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

func (s *PostgresStorage) PingDBServer(ctx context.Context) (int, error) {
	err := s.db.PingContext(ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
