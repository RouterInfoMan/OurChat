package repositories

import (
	"database/sql"
	"errors"
	"github.com/RouterInfoMan/OurChat/internal/models"
)

// DatabaseRepo handles database operations
type DatabaseRepo struct {
	db *sql.DB
}

// NewDatabaseRepo creates a new DatabaseRepo instance
func NewDatabaseRepo(dbFile string) (*DatabaseRepo, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	
	// Create users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}
	
	return &DatabaseRepo{db: db}, nil
}

// Close closes the database connection
func (dr *DatabaseRepo) Close() error {
	return dr.db.Close()
}

// RegisterUser registers a new user in the database
func (dr *DatabaseRepo) RegisterUser(username, email, password string) (int64, error) {
	// Check if username already exists
	var exists int
	err := dr.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists > 0 {
		return 0, ErrUsernameTaken
	}
	
	// Check if email already exists
	err = dr.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists > 0 {
		return 0, ErrEmailTaken
	}
	
	// Insert new user
	result, err := dr.db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, password,
	)
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// GetUserByUsername fetches a user by their username
func (dr *DatabaseRepo) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := dr.db.QueryRow(
		"SELECT id, username, email, password, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

// Error definitions
// The `var` block you provided is defining error variables using the `errors.New` function. Each error
// variable is created with a specific error message associated with it. These error variables are used
// to represent different error conditions that can occur during the execution of the code.
var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrEmailTaken    = errors.New("email already taken")
	ErrUserNotFound  = errors.New("user not found")
)