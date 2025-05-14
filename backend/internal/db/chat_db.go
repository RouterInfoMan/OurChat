package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"OurChat/internal/models"
)

// CreateChat creates a new chat with the specified type and name
func (db *DB) CreateChat(chatType, name string) (int64, error) {
	query := `
	INSERT INTO chats (type, name, created_at, updated_at)
	VALUES (?, ?, ?, ?)`

	now := time.Now()
	result, err := db.Exec(query, chatType, name, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	log.Println("Chat created successfully")
	return id, nil
}

// AddUserToChat adds a user to a chat with the specified role
func (db *DB) AddUserToChat(userID, chatID int, role string) error {
	query := `
	INSERT INTO chat_members (user_id, chat_id, role, joined_at)
	VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, userID, chatID, role, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add user to chat: %w", err)
	}

	log.Println("User added to chat successfully")
	return nil
}

// GetChatsForUser retrieves all chats that a user is a member of
func (db *DB) GetChatsForUser(userID int) ([]models.Chat, error) {
	query := `
	SELECT c.id, c.type, c.name, c.created_at, c.updated_at, c.is_active
	FROM chats c
	JOIN chat_members cm ON c.id = cm.chat_id
	WHERE cm.user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.Name, &chat.CreatedAt, &chat.UpdatedAt, &chat.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	log.Println("Chats retrieved successfully")
	return chats, nil
}

// GetChatByID retrieves a chat by its ID
func (db *DB) GetChatByID(chatID int) (*models.Chat, error) {
	chat := &models.Chat{}
	query := `SELECT id, type, name, created_at, updated_at, is_active
	          FROM chats WHERE id = ?`

	err := db.QueryRow(query, chatID).Scan(
		&chat.ID, &chat.Type, &chat.Name,
		&chat.CreatedAt, &chat.UpdatedAt, &chat.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	log.Println("Chat retrieved successfully")
	return chat, nil
}

// GetUsersByChatID retrieves all users who are members of a specific chat
func (db *DB) GetUsersByChatID(chatID int) ([]models.User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.status
	FROM users u
	JOIN chat_members cm ON u.id = cm.user_id
	WHERE cm.chat_id = ?`

	rows, err := db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	log.Println("Chat users retrieved successfully")
	return users, nil
}

// GetDirectChatBetweenUsers finds or creates a direct chat between two users
func (db *DB) GetDirectChatBetweenUsers(userID1, userID2 int) (int, error) {
	// First try to find an existing direct chat between these users
	query := `
	SELECT c.id FROM chats c
	JOIN chat_members cm1 ON c.id = cm1.chat_id AND cm1.user_id = ?
	JOIN chat_members cm2 ON c.id = cm2.chat_id AND cm2.user_id = ?
	WHERE c.type = 'direct'
	LIMIT 1`

	var chatID int
	err := db.QueryRow(query, userID1, userID2).Scan(&chatID)
	if err == nil {
		// Found existing chat
		return chatID, nil
	}

	if err != sql.ErrNoRows {
		// Unexpected error
		return 0, fmt.Errorf("failed to find direct chat: %w", err)
	}

	// No existing chat found, create a new one
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the chat
	chatIDInt64, err := db.CreateChat("direct", "")
	if err != nil {
		return 0, fmt.Errorf("failed to create direct chat: %w", err)
	}

	chatID = int(chatIDInt64)

	// Add both users to the chat
	if err := db.AddUserToChat(userID1, chatID, "member"); err != nil {
		return 0, fmt.Errorf("failed to add user1 to chat: %w", err)
	}

	if err := db.AddUserToChat(userID2, chatID, "member"); err != nil {
		return 0, fmt.Errorf("failed to add user2 to chat: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Created new direct chat %d between users %d and %d", chatID, userID1, userID2)
	return chatID, nil
}
