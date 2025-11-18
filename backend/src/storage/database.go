package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Initialize creates a new database connection and runs migrations
func Initialize(dbPath string) (*sql.DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Run migrations
	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	// Ensure admin user exists
	if err := EnsureAdminUser(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
