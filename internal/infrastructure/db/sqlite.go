package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(homeDir, ".govault.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id         TEXT PRIMARY KEY,
		username   TEXT NOT NULL UNIQUE,
		password   TEXT NOT NULL,
		salt       TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS credentials (
		id         TEXT PRIMARY KEY,
		user_id    TEXT NOT NULL,
		site       TEXT NOT NULL,
		username   TEXT NOT NULL,
		password   TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(query)
	return err
}
