package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"OurChat/internal/models"
)

// CreateUser creates a new user record with a random JWT key
func (db *DB) CreateUser(username, email, password string) error {
	// Generate a random JWT key
	jwtKey, err := GenerateJWTKey()
	if err != nil {
		return fmt.Errorf("failed to generate JWT key: %w", err)
	}

	query := `
	INSERT INTO users (username, email, password, jwt_key, created_at)
	VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, username, email, password, jwtKey, time.Now())

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Println("User created successfully")
	return nil
}

// GetUserByUsername retrieves a user by their username
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
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
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	log.Println("User retrieved successfully")
	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, jwt_key, created_at, last_login, status
	          FROM users WHERE email = ?`

	var lastLogin sql.NullTime

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.JWTKey, &user.CreatedAt, &lastLogin, &user.Status,
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

	log.Println("User retrieved successfully")
	return user, nil
}

// UpdateUserStatus updates a user's online/offline status
func (db *DB) UpdateUserStatus(userID int, status string) error {
	query := `UPDATE users SET status = ? WHERE id = ?`
	_, err := db.Exec(query, status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	log.Println("User status updated successfully")
	return nil
}

// UpdateUserProfile updates a user's profile information
func (db *DB) UpdateUserProfile(userID int, updates map[string]interface{}) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Process each update field
	for field, value := range updates {
		var query string
		var args []interface{}

		switch field {
		case "email":
			// Check if email is already in use by another user
			var count int
			checkQuery := "SELECT COUNT(*) FROM users WHERE email = ? AND id != ?"
			err := db.QueryRow(checkQuery, value, userID).Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check email uniqueness: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("email is already in use")
			}

			query = "UPDATE users SET email = ? WHERE id = ?"
			args = []interface{}{value, userID}

		case "status":
			query = "UPDATE users SET status = ? WHERE id = ?"
			args = []interface{}{value, userID}

		default:
			// Skip unsupported fields
			continue
		}

		// Execute the update
		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to update %s: %w", field, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUsersByIDs retrieves basic information for a list of user IDs
func (db *DB) GetUsersByIDs(userIDs []int) (map[int]models.UserBasic, error) {
	if len(userIDs) == 0 {
		return make(map[int]models.UserBasic), nil
	}

	// Create placeholders for SQL query (?, ?, ?, etc.)
	placeholders := strings.Repeat("?,", len(userIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf(`
    SELECT id, username, status, profile_picture_url
    FROM users
    WHERE id IN (%s)`, placeholders)

	// Convert userIDs to []interface{} for db.Query
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := make(map[int]models.UserBasic)
	for rows.Next() {
		var user models.UserBasic
		if err := rows.Scan(&user.ID, &user.Username, &user.Status, &user.ProfilePictureURL); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users[user.ID] = user
	}

	return users, nil
}

// SearchUsersByName searches for users by partial username match
func (db *DB) SearchUsersByName(searchTerm string, limit int) ([]models.UserBasic, error) {
	if searchTerm == "" {
		return make([]models.UserBasic, 0), nil
	}

	// Validate limit
	if limit <= 0 || limit > 50 {
		limit = 20 // Default limit
	}

	// Add wildcards for SQL LIKE search
	searchPattern := "%" + searchTerm + "%"

	query := `
    SELECT id, username, status, profile_picture_url
    FROM users
    WHERE username LIKE ?
    ORDER BY
        CASE
            WHEN username = ? THEN 1
            WHEN username LIKE ? THEN 2
            ELSE 3
        END,
        username ASC
    LIMIT ?`

	// Parameters for the query:
	// 1. searchPattern for the WHERE clause
	// 2. searchTerm for exact match priority
	// 3. searchTerm + "%" for starts-with match priority
	startsWithPattern := searchTerm + "%"

	rows, err := db.Query(query, searchPattern, searchTerm, startsWithPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	users := make([]models.UserBasic, 0)
	for rows.Next() {
		var user models.UserBasic
		if err := rows.Scan(&user.ID, &user.Username, &user.Status, &user.ProfilePictureURL); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	log.Printf("Found %d users matching '%s'", len(users), searchTerm)
	return users, nil
}
