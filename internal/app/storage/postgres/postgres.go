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
	urlID, err := getURLIDFromQueryRow(row)
	if err != nil {
		return fmt.Errorf("error while save url processing in postgres: %w", err)
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

func (s *PostgresStorage) SaveBatchURL(ctx context.Context, userID entity.UserID, batch model.Batch) (model.Batch, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create save batch url transaction in postgres: %w", err)
	}
	defer tx.Rollback()

	saveURLQuery := `INSERT INTO url(short_url, url) VALUES($1, $2) RETURNING id`
	stmtURL, err := tx.PrepareContext(ctx, saveURLQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare save batch URL query in postgres: %w", err)
	}
	defer stmtURL.Close()

	saveUserURLQuery := `INSERT INTO users_url VALUES($1, $2)`
	stmtUserURL, err := tx.PrepareContext(ctx, saveUserURLQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare save batch user URL query in postgres: %w", err)
	}
	defer stmtURL.Close()

	for _, obj := range batch {
		row := stmtURL.QueryRowContext(ctx, obj.ShortURL, obj.InputURL)
		urlID, err := getURLIDFromQueryRow(row)
		if err != nil {
			return nil, fmt.Errorf("error while save batch url processing in postgres: %w", err)
		}

		_, err = stmtUserURL.ExecContext(ctx, userID.String(), urlID)
		if err != nil {
			return nil, fmt.Errorf("failed to save batch users url object to postgres: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("unable to commit save batch url transaction in postgres: %w", err)
	}

	return batch, nil
}

func (s *PostgresStorage) GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error) {
	query := `SELECT u.url FROM users_url AS uu
					JOIN url AS u
						ON uu.url_id = u.id
				WHERE uu.users_id = @userID AND u.short_url = @shortUrl`
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

func (s *PostgresStorage) GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error) {
	query := `SELECT u.url, u.short_url FROM users_url AS uu
				JOIN url AS u
					ON uu.url_id = u.id
			  WHERE uu.users_id = @userID`

	args := pgx.NamedArgs{
		"userID":   userID.String(),
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

func (s *PostgresStorage) AddUser(ctx context.Context, userID entity.UserID) error {
	query := `INSERT INTO users VALUES(@usersId)`
	args := pgx.NamedArgs{
		"usersId": userID.String(),
	}

	_, err := s.db.ExecContext(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("error while adding user to postgres: %w", api.ErrUserAlreadyExists)
		}

		return fmt.Errorf("unable to add user to postgres: %w", err)
	}

	return nil
}

func (s *PostgresStorage) AuthUser(ctx context.Context, userID entity.UserID) (entity.UserID, error) {
	query := `SELECT * FROM users WHERE id=$1`
	row := s.db.QueryRowContext(ctx, query, userID.String())
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

func getURLIDFromQueryRow(row *sql.Row) (int, error) {
	if row == nil {
		return 0, fmt.Errorf("row is nil")
	}

	if row.Err() != nil {
		return 0, fmt.Errorf("request execution row error: %w", row.Err())
	}

	var urlID int
	err := row.Scan(&urlID)
	if err != nil {
		return 0, fmt.Errorf("error while scan row: %w", err)
	}

	return urlID, nil
}
