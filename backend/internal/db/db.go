package db

// Core DB functions

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps a sql.DB connection and provides access to all database operations
type DB struct {
	*sql.DB
}

// NewDB initializes a new database connection and loads the schema if needed
func NewDB(dbPath string) (*DB, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create DB wrapper
	database := &DB{db}

	// Load schema if tables don't exist
	if err := database.LoadSchemaIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	return database, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
