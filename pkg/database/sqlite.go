package database

import (
	"database/sql"
	"fmt"
	"os"

	"forum/config"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initialize sqlite3 database and checks established connection
func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DbDriver, cfg.DbNameAndPath)
	if err != nil {
		return nil, err
	}
	// checks connection to db
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateTables executes all tables for forum
func CreateTables(db *sql.DB) error {
	migrationData, err := os.ReadFile("./migrations/up.sql")
	if err != nil {
		return fmt.Errorf("create tables: read file: %w", err)
	}
	if _, err = db.Exec(string(migrationData)); err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
