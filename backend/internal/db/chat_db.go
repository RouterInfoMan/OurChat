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
	chats = make([]models.Chat, 0)
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.Name, &chat.CreatedAt, &chat.UpdatedAt, &chat.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	for chat := range chats {
		// Put the name of the other user in the chat name for direct chats
		if chats[chat].Type == "direct" {
			otherUserQuery := `
			SELECT u.username FROM chat_members cm
			JOIN users u ON cm.user_id = u.id
			WHERE cm.chat_id = ? AND cm.user_id != ?`
			var otherUsername string
			err := db.QueryRow(otherUserQuery, chats[chat].ID, userID).Scan(&otherUsername)
			if err != nil {
				return nil, fmt.Errorf("failed to get other user for direct chat: %w", err)
			}
			chats[chat].Name = otherUsername
		}
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

// GetChatMembers returns all users who are members of a specific chat
func (db *DB) GetChatMembers(chatID int) ([]models.ChatMember, error) {
	query := `
    SELECT cm.id, cm.user_id, cm.chat_id, cm.role, cm.joined_at, cm.last_read_at,
           u.username, u.status, u.profile_picture_url
    FROM chat_members cm
    JOIN users u ON cm.user_id = u.id
    WHERE cm.chat_id = ?
    ORDER BY cm.role = 'admin' DESC, u.username ASC`

	rows, err := db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat members: %w", err)
	}
	defer rows.Close()

	var members []models.ChatMember
	for rows.Next() {
		var member models.ChatMember
		var lastReadAt sql.NullTime

		err := rows.Scan(
			&member.ID,
			&member.UserID,
			&member.ChatID,
			&member.Role,
			&member.JoinedAt,
			&lastReadAt,
			&member.Username,
			&member.Status,
			&member.ProfilePictureURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat member: %w", err)
		}

		if lastReadAt.Valid {
			member.LastReadAt = &lastReadAt.Time
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat members: %w", err)
	}

	return members, nil
}
