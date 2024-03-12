package migration

import (
	"database/sql"
	"fmt"
)

func InitDBTables(db *sql.DB) error {
	if err := createURLTable(db); err != nil {
		return fmt.Errorf("error while creating url table: %w", err)
	}

	return nil
}

func createURLTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			short_url TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_short_url ON url(short_url);
	`)
	
	if err != nil {
		return err
	}

	return nil
}
