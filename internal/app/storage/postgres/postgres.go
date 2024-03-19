package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/avGenie/url-shortener/internal/app/storage/postgres/migration"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	saveQuery      = `INSERT INTO url(short_url, url) VALUES(@shortUrl, @url)`
	saveBatchQuery = `INSERT INTO url(short_url, url) VALUES($1, $2)`
	getQuery       = `SELECT url FROM url WHERE short_url=@shortUrl`
)

type PostgresStorage struct {
	model.Storage

	db *sql.DB
}

func NewPostgresStorage(dbStorageConnect string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dbStorageConnect)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql connect: %w", err)
	}

	err = migration.InitDBTables(db)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql table initialization, %w", err)
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

func (s *PostgresStorage) PingServer(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("couldn't ping postgres server: %w", err)
	}

	return nil
}

func (s *PostgresStorage) SaveURL(ctx context.Context, key, value entity.URL) error {
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
		"url":      value.String(),
	}

	_, err := s.db.ExecContext(ctx, saveQuery, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("error while save url to postgres: %w", api.ErrURLAlreadyExists)
        }
		
		return fmt.Errorf("unable to insert row to postgres: %w", err)
	}

	return nil
}

func (s *PostgresStorage) SaveBatchURL(ctx context.Context, batch model.Batch) (model.Batch, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction in postgres: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, saveBatchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query in postgres: %w", err)
	}
	defer stmt.Close()

	for _, obj := range batch {
		_, err = stmt.ExecContext(ctx, obj.ShortURL, obj.InputURL)
		if err != nil {
			return nil, fmt.Errorf("failed to write batch object to postgres: %w", err)
		}
	}
	tx.Commit()

	return batch, nil
}

func (s *PostgresStorage) GetURL(ctx context.Context, key entity.URL) (*entity.URL, error) {
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
	}

	var dbURL string
	row := s.db.QueryRowContext(ctx, getQuery, args)
	if row == nil {
		return nil, fmt.Errorf("error while postgres request execution")
	}

	if row.Err() != nil {
		return nil, fmt.Errorf("error while postgres request execution: %w", row.Err())
	}

	err := row.Scan(&dbURL)
	if err != nil {
		return nil, fmt.Errorf("error while processing response row in postgres: %w", err)
	}

	url, err := entity.NewURL(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error while creating url in postgres: %w", err)
	}

	return url, nil
}
