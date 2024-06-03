// Package postgres contains implementation of postgres storage
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/models"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

const (
	migrationDB     = "postgres"
	migrationFolder = "migrations"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

// PostgresStorage Postgres storage object
type PostgresStorage struct {
	model.Storage

	db *sql.DB
}

// NewPostgresStorage Creates postgres storage object
func NewPostgresStorage(dbStorageConnect string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dbStorageConnect)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql connect: %w", err)
	}

	err = migration(db)
	if err != nil {
		return nil, fmt.Errorf("error while postgresql migration: %w", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

// Close Closes connection to postgres DB
func (s *PostgresStorage) Close() {
	err := s.db.Close()
	if err != nil {
		zap.L().Error("error while closing postgres storage", zap.Error(err))
	}
}

// PingServer Pings to postgres DB server
func (s *PostgresStorage) PingServer(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("couldn't ping postgres server: %w", err)
	}

	return nil
}

// SaveURL Saves user URL to postgres DB
func (s *PostgresStorage) SaveURL(ctx context.Context, userID entity.UserID, key, value entity.URL) error {
	query := `INSERT INTO url(short_url, url, user_id) VALUES(@shortUrl, @url, @userID)`
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
		"url":      value.String(),
		"userID":   userID.String(),
	}

	_, err := s.db.ExecContext(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("error while save url to postgres: %w", api.ErrURLAlreadyExists)
		}

		return fmt.Errorf("unable to insert row to postgres: %w", err)
	}

	return nil
}

// SaveBatchURL Saves batch of user URLs to postgres DB
func (s *PostgresStorage) SaveBatchURL(ctx context.Context, userID entity.UserID, batch model.Batch) (model.Batch, error) {
	query := `INSERT INTO url(short_url, url, user_id) VALUES($1, $2, $3)`
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction in postgres: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query in postgres: %w", err)
	}
	defer stmt.Close()

	for _, obj := range batch {
		_, err = stmt.ExecContext(ctx, obj.ShortURL, obj.InputURL, userID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to write batch object to postgres: %w", err)
		}
	}
	tx.Commit()

	return batch, nil
}

// GetURL Returns URL user from postgres DB
func (s *PostgresStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	if !userID.IsValid() {
		return s.getURL(ctx, key)
	}

	return s.getUserURL(ctx, userID, key)
}

// GetAllURLByUserID Returns all user URLs from postgres DB
func (s *PostgresStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get all urls transaction in postgres: %w", err)
	}
	defer tx.Rollback()

	query := `SELECT url, short_url, deleted FROM url WHERE user_id = @userID`
	args := pgx.NamedArgs{
		"userID": userID.String(),
	}

	rows, err := s.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("error in postgres request execution while getting all urls by user id: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in postgres requested rows while getting all urls by user id: %w", rows.Err())
	}

	isDeleted := false
	var urlsBatch models.AllUrlsBatch
	for rows.Next() {
		var url models.AllUrlsResponse
		var deleted bool
		err := rows.Scan(&url.OriginalURL, &url.ShortURL, &deleted)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, api.ErrShortURLNotFound
			}
			return nil, fmt.Errorf("error while processing response row in postgres: %w", err)
		}

		if deleted {
			isDeleted = true
			err = deleteURL(s.db, ctx, userID.String(), url.ShortURL)
			if err != nil {
				zap.L().Error(
					"unable to delete url while getting all urls from postgres",
					zap.Error(err),
					zap.String("short_url", url.ShortURL))
			}

			continue
		}

		urlsBatch = append(urlsBatch, url)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction while getting all urls in postgres: %w", err)
	}

	if len(urlsBatch) == 0 && isDeleted {
		return nil, api.ErrAllURLsDeleted
	}

	return urlsBatch, nil
}

// DeleteBatchURL Delete user URL from postgres DB
func (s *PostgresStorage) DeleteBatchURL(ctx context.Context, urls entity.DeletedURLBatch) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete batch url transaction in postgres: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE url SET deleted = true WHERE user_id=$1 AND short_url=$2`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query while deleting urls in postgres: %w", err)
	}
	defer stmt.Close()

	for _, url := range urls {
		_, err = stmt.ExecContext(ctx, url.UserID, url.ShortURL)
		if err != nil {
			return fmt.Errorf("failed to update deleted url in postgres: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit deleted batch url transaction in postgres: %w", err)
	}

	return nil
}

func (s *PostgresStorage) getURL(ctx context.Context, key entity.URL) (*entity.URL, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create get user url transaction in postgres: %w", err)
	}
	query := `SELECT id, url, deleted FROM url WHERE short_url=@shortUrl`
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
	}

	row := s.db.QueryRowContext(ctx, query, args)
	if row == nil {
		return nil, fmt.Errorf("error while postgres request execution")
	}

	if row.Err() != nil {
		return nil, fmt.Errorf("error while postgres request execution: %w", row.Err())
	}

	var id int
	var dbURL string
	var deleted bool
	err = row.Scan(&id, &dbURL, &deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error while processing response row in postgres: %w", err)
	}

	if deleted {
		query = `DELETE FROM url WHERE id=$1`
		_, err := s.db.ExecContext(ctx, query, id)
		if err != nil {
			zap.L().Error(
				"unable to delete url while getting from postgres",
				zap.Error(err),
				zap.String("short_url", key.String()))
		}

		return nil, api.ErrAllURLsDeleted
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction while getting user url in postgres: %w", err)
	}

	url, err := entity.NewURL(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error while creating url in postgres: %w", err)
	}

	return url, nil
}

func (s *PostgresStorage) getUserURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create get user url transaction in postgres: %w", err)
	}
	query := `SELECT url, deleted FROM url WHERE user_id = @userID AND short_url = @shortUrl`
	args := pgx.NamedArgs{
		"userID":   userID.String(),
		"shortUrl": key.String(),
	}

	row := s.db.QueryRowContext(ctx, query, args)
	if row == nil {
		return nil, fmt.Errorf("error in postgres request execution while getting url")
	}

	if row.Err() != nil {
		return nil, fmt.Errorf("error in postgres request execution while getting url: %w", row.Err())
	}

	var dbURL string
	var deleted bool
	err = row.Scan(&dbURL, &deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error in postgres processing response row while getting url: %w", err)
	}

	if deleted {
		err = deleteURL(s.db, ctx, userID.String(), key.String())
		if err != nil {
			zap.L().Error(
				"unable to delete url while getting from postgres",
				zap.Error(err),
				zap.String("short_url", key.String()))
		}

		return nil, api.ErrAllURLsDeleted
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction while getting user url in postgres: %w", err)
	}

	url, err := entity.NewURL(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error in postgres creating url while getting url: %w", err)
	}

	return url, nil
}

func migration(db *sql.DB) error {
	goose.SetBaseFS(migrationFs)

	if err := goose.SetDialect(migrationDB); err != nil {
		return err
	}

	if err := goose.Up(db, migrationFolder); err != nil {
		return err
	}

	return nil
}

func deleteURL(db *sql.DB, ctx context.Context, userID string, shortURL string) error {
	query := `DELETE FROM url WHERE user_id=$1 AND short_url=$2`
	_, err := db.ExecContext(ctx, query, userID, shortURL)
	if err != nil {
		return fmt.Errorf("unable to delete row from postgres: %w", err)
	}

	return nil
}
