package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// LoadSchemaIfNeeded loads the SQL schema file if the tables don't exist yet
func (db *DB) LoadSchemaIfNeeded() error {
	// Check if users table exists (as a proxy for all tables)
	var count int
	err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if tables exist: %w", err)
	}

	// If tables already exist, return early
	if count > 0 {
		log.Println("Database schema already exists")
		return nil
	}

	// Load schema file
	schemaPath := filepath.Join("internal", "db", "migrations", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema SQL
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}

	log.Println("Database schema loaded successfully")
	return nil
}
