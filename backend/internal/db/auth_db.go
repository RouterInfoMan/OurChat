package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"OurChat/internal/models"
)

// GenerateJWTKey creates a cryptographically secure random key for JWT signing
func GenerateJWTKey() (string, error) {
	// Generate 32 bytes of random data for the key
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 for storage
	key := base64.StdEncoding.EncodeToString(bytes)
	return key, nil
}

// GetUserForAuth retrieves a user by username for authentication purposes
// This method is specifically designed for login/auth rather than general user retrieval
func (db *DB) GetUserForAuth(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, jwt_key, created_at, last_login, status
	          FROM users WHERE username = ?`

	var lastLogin sql.NullTime

	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.JWTKey, &user.CreatedAt, &lastLogin, &user.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("authentication failed: user not found")
		}
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetUserByID retrieves a user by ID (useful for JWT middleware)
func (db *DB) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, jwt_key, profile_picture_url, created_at, last_login, status
	          FROM users WHERE id = ?`

	var lastLogin sql.NullTime
	var profilePictureURL sql.NullString

	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.JWTKey, &profilePictureURL, &user.CreatedAt, &lastLogin, &user.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	if profilePictureURL.Valid {
		user.ProfilePictureURL = &profilePictureURL.String
	}

	return user, nil
}

// UpdateJWTKey generates and updates the JWT key for a specific user
// This can be used to invalidate all existing tokens when needed
func (db *DB) UpdateJWTKey(userID int) (string, error) {
	// Generate a new random JWT key
	newKey, err := GenerateJWTKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate new JWT key: %w", err)
	}

	// Update the key in the database
	query := `UPDATE users SET jwt_key = ? WHERE id = ?`
	_, err = db.Exec(query, newKey, userID)
	if err != nil {
		return "", fmt.Errorf("failed to update JWT key: %w", err)
	}

	return newKey, nil
}

// UpdateLastLogin updates the last_login timestamp and online status for a user
func (db *DB) UpdateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = ?, status = 'online' WHERE id = ?`
	_, err := db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// LogoutUser updates the user's status to 'offline'
func (db *DB) LogoutUser(userID int) error {
	query := `UPDATE users SET status = 'offline' WHERE id = ?`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// GetUserJWTKey retrieves only the JWT key for a user by ID
// This is useful for token validation in middleware
func (db *DB) GetUserJWTKey(userID int) (string, error) {
	var jwtKey string
	query := `SELECT jwt_key FROM users WHERE id = ?`

	err := db.QueryRow(query, userID).Scan(&jwtKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found: %w", err)
		}
		return "", fmt.Errorf("failed to get JWT key: %w", err)
	}

	return jwtKey, nil
}

// RequestPasswordReset initiates a password reset for a user
func (db *DB) RequestPasswordReset(email string) (*models.User, error) {
	// Check if the user exists
	user := &models.User{}
	query := `SELECT id, username, email, jwt_key FROM users WHERE email = ?`

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.JWTKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with this email")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ResetPassword changes a user's password and rotates their JWT key to invalidate existing tokens
func (db *DB) ResetPassword(userID int, newPassword string) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update the password
	_, err = tx.Exec("UPDATE users SET password = ? WHERE id = ?", newPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Generate a new JWT key to invalidate existing tokens
	newJWTKey, err := GenerateJWTKey()
	if err != nil {
		return fmt.Errorf("failed to generate new JWT key: %w", err)
	}

	// Update the JWT key
	_, err = tx.Exec("UPDATE users SET jwt_key = ? WHERE id = ?", newJWTKey, userID)
	if err != nil {
		return fmt.Errorf("failed to update JWT key: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
