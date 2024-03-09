package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	entity.Storage

	db *sql.DB
}

func NewPostgresStorage(dbStorageConnect string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dbStorageConnect)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql connect: %w", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Close() entity.Response {
	err := s.db.Close()
	if err != nil {
		outErr := fmt.Errorf("couldn'r closed postgres db: %w", err)
		return entity.ErrorResponse(outErr)
	}

	return entity.OKResponse()
}

func (s *PostgresStorage) PingServer(ctx context.Context) entity.Response {
	err := s.db.PingContext(ctx)
	if err != nil {
		outErr := fmt.Errorf("couldn'r ping postgres server: %w", err)
		return entity.ErrorResponse(outErr)
	}

	return entity.OKResponse()
}

func (s *PostgresStorage) AddURL(ctx context.Context, key, value entity.URL) entity.Response {
	return entity.OKResponse()
}

func (s *PostgresStorage) GetURL(ctx context.Context, key entity.URL) entity.URLResponse {
	return entity.OKURLResponse(entity.URL{})
}
