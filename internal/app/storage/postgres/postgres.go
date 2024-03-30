package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/avGenie/url-shortener/internal/app/entity"
	api "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	addUserQuery   = `INSERT INTO users VALUES(@usersId)`
	saveBatchQuery = `INSERT INTO url(short_url, url) VALUES($1, $2)`
	getQuery       = `SELECT url FROM url WHERE short_url=@shortUrl`
	getUser        = `SELECT * FROM users WHERE id=$1`

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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin save url transaction in postgres: %w", err)
	}

	defer tx.Rollback()

	saveURLQuery := `INSERT INTO url(short_url, url) VALUES(@shortUrl, @url) RETURNING id`
	argsURL := pgx.NamedArgs{
		"shortUrl": key.String(),
		"url":      value.String(),
	}

	row := s.db.QueryRowContext(ctx, saveURLQuery, argsURL)
	if row == nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			fmt.Println("QueryRowContext pgx")
			return fmt.Errorf("unable to save url in postgres: %w", api.ErrURLAlreadyExists)
		}
		fmt.Println("QueryRowContext")

		return fmt.Errorf("unable to save url in postgres: %w", err)
	}

	if row.Err() != nil {
		return fmt.Errorf("error while postgres save url request execution: %w", row.Err())
	}

	var urlID int
	err = row.Scan(&urlID)
	if err != nil {
		return fmt.Errorf("error while scan url id in postgres: %w", err)
	}

	saveURLUserQuery := `INSERT INTO users_url VALUES(@userID, @urlID)`
	argsURLUser := pgx.NamedArgs{
		"userID": userID.String(),
		"urlID":  urlID,
	}

	_, err = s.db.ExecContext(ctx, saveURLUserQuery, argsURLUser)
	if err != nil {
		return fmt.Errorf("unable to save users and url in postgres: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit save url transaction in postgres: %w", err)
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

func (s *PostgresStorage) AddUser(ctx context.Context, userID entity.UserID) error {
	args := pgx.NamedArgs{
		"usersId": userID.String(),
	}

	_, err := s.db.ExecContext(ctx, addUserQuery, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("error while adding user to postgres: %w", api.ErrURLAlreadyExists)
		}

		return fmt.Errorf("unable to add user to postgres: %w", err)
	}

	return nil
}

func (s *PostgresStorage) AuthUser(ctx context.Context, userID entity.UserID) (entity.UserID, error) {
	row := s.db.QueryRowContext(ctx, getUser, userID.String())
	if row == nil {
		return "", fmt.Errorf("error while postgres prepare row")
	}

	if row.Err() != nil {
		return "", fmt.Errorf("error while postgres request execution: %w", row.Err())
	}
	var id string
	err := row.Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		if err == sql.ErrNoRows {
			return "", api.ErrUserIDNotFound
		}
		return "", fmt.Errorf("error while processing response row in postgres: %w", err)
	}

	return userID, nil
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
