package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"OurChat/internal/models"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return database, nil
}

func (db *DB) createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	// TODO Change receiver_id to be a foreign key to chats table
	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id),
		FOREIGN KEY (receiver_id) REFERENCES users(id)
	);`

	_, err := db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Println("Users table created successfully")

	_, err = db.Exec(messagesTable)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}
	log.Println("Messages table created successfully")

	return nil
}
func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) CreateUser(username, email, password string) error {
	query := `
	INSERT INTO users (username, email, password, created_at)
	VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, username, email, password, time.Now())

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Println("User created successfully")
	return nil
}

func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE username = ?`

	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	log.Println("User retrieved successfully")
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	log.Println("User retrieved successfully")
	return user, nil
}
func (db *DB) GetMessagesByUserID(userID int) ([]models.Message, error) {
	query := `SELECT id, sender_id, receiver_id, content, created_at FROM messages WHERE sender_id = ? OR receiver_id = ?`
	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}
	log.Println("Messages retrieved successfully")
	return messages, nil
}

func (db *DB) GetMessagesByUserToUser(userID1, userID2 int) ([]models.Message, error) {
	query := `SELECT id, sender_id, receiver_id, content, created_at FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)`
	rows, err := db.Query(query, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}
	log.Println("Messages retrieved successfully")
	return messages, nil
}

func (db *DB) CreateMessage(senderID, receiverID int, content string) error {
	query := `
	INSERT INTO messages (sender_id, receiver_id, content, created_at)
	VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, senderID, receiverID, content, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	log.Println("Message created successfully")
	return nil
}
