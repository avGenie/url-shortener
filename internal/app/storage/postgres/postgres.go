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
)

const (
	migrationDB     = "postgres"
	migrationFolder = "migrations"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

type PostgresStorage struct {
	model.Storage

	db *sql.DB
}

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

func (s *PostgresStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	if !userID.IsValid() {
		return s.getURL(ctx, key)
	}

	return s.getUserURL(ctx, userID, key)
}

func (s *PostgresStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	query := `SELECT url, short_url FROM url WHERE user_id = @userID`
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

	var urlsBatch models.AllUrlsBatch
	for rows.Next() {
		var allURLs models.AllUrlsResponse
		err := rows.Scan(&allURLs.OriginalURL, &allURLs.ShortURL)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, api.ErrShortURLNotFound
			}
			return nil, fmt.Errorf("error while processing response row in postgres: %w", err)
		}

		urlsBatch = append(urlsBatch, allURLs)
	}

	return urlsBatch, nil
}

func (s *PostgresStorage) getURL(ctx context.Context, key entity.URL) (*entity.URL, error) {
	query := `SELECT url FROM url WHERE short_url=@shortUrl`
	args := pgx.NamedArgs{
		"shortUrl": key.String(),
	}

	var dbURL string
	row := s.db.QueryRowContext(ctx, query, args)
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

func (s *PostgresStorage) getUserURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	query := `SELECT url FROM url WHERE user_id = @userID AND short_url = @shortUrl`
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
	err := row.Scan(&dbURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error in postgres processing response row while getting url: %w", err)
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
