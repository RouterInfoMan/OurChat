package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"OurChat/internal/models"
)

func (db *DB) CreateMessage(senderID, chatID int, content, messageType string, mediaFileID *int) (int64, error) {
	query := `
	INSERT INTO messages (sender_id, chat_id, content, message_type, media_file_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, senderID, chatID, content, messageType, mediaFileID, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	// Update the chat's updated_at timestamp
	updateQuery := `UPDATE chats SET updated_at = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, time.Now(), chatID)
	if err != nil {
		log.Printf("Warning: failed to update chat timestamp: %v", err)
	}

	log.Println("Message created successfully")
	return id, nil
}

// GetMessageByID retrieves a specific message by its ID
func (db *DB) GetMessageByIDWithMedia(messageID int) (*models.Message, error) {
	message := &models.Message{}
	query := `
	SELECT m.id, m.sender_id, m.chat_id, m.content, m.message_type, m.media_file_id, m.created_at, m.is_read,
	       mf.id, mf.filename, mf.original_filename, mf.file_size, mf.mime_type, mf.uploaded_at
	FROM messages m
	LEFT JOIN media_files mf ON m.media_file_id = mf.id
	WHERE m.id = ?`

	var mediaFileID, mediaID sql.NullInt64
	var mediaFilename, mediaOriginalFilename, mediaMimeType sql.NullString
	var mediaFileSize sql.NullInt64
	var mediaUploadedAt sql.NullTime

	err := db.QueryRow(query, messageID).Scan(
		&message.ID, &message.SenderID, &message.ChatID, &message.Content,
		&message.MessageType, &mediaFileID, &message.CreatedAt, &message.IsRead,
		&mediaID, &mediaFilename, &mediaOriginalFilename, &mediaFileSize,
		&mediaMimeType, &mediaUploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Set media file ID if it exists
	if mediaFileID.Valid {
		id := int(mediaFileID.Int64)
		message.MediaFileID = &id
	}

	// Populate media file information if it exists
	if mediaID.Valid {
		message.MediaFile = &models.MediaFile{
			ID:               int(mediaID.Int64),
			Filename:         mediaFilename.String,
			OriginalFilename: mediaOriginalFilename.String,
			FileSize:         mediaFileSize.Int64,
			MimeType:         mediaMimeType.String,
			UploadedAt:       mediaUploadedAt.Time,
			URL:              fmt.Sprintf("/api/media/files/%s", mediaFilename.String),
		}
	}

	return message, nil
}

// GetMessagesByChatID retrieves messages from a chat with pagination
func (db *DB) GetMessagesByChatIDWithMedia(chatID int, limit, offset int) ([]models.Message, error) {
	query := `
	SELECT m.id, m.sender_id, m.chat_id, m.content, m.message_type, m.media_file_id, m.created_at, m.is_read,
	       mf.id, mf.filename, mf.original_filename, mf.file_size, mf.mime_type, mf.uploaded_at
	FROM messages m
	LEFT JOIN media_files mf ON m.media_file_id = mf.id
	WHERE m.chat_id = ?
	ORDER BY m.created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := db.Query(query, chatID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		var mediaFileID, mediaID sql.NullInt64
		var mediaFilename, mediaOriginalFilename, mediaMimeType sql.NullString
		var mediaFileSize sql.NullInt64
		var mediaUploadedAt sql.NullTime

		err := rows.Scan(
			&message.ID, &message.SenderID, &message.ChatID, &message.Content,
			&message.MessageType, &mediaFileID, &message.CreatedAt, &message.IsRead,
			&mediaID, &mediaFilename, &mediaOriginalFilename, &mediaFileSize,
			&mediaMimeType, &mediaUploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Set media file ID if it exists
		if mediaFileID.Valid {
			id := int(mediaFileID.Int64)
			message.MediaFileID = &id
		}

		// Populate media file information if it exists
		if mediaID.Valid {
			message.MediaFile = &models.MediaFile{
				ID:               int(mediaID.Int64),
				Filename:         mediaFilename.String,
				OriginalFilename: mediaOriginalFilename.String,
				FileSize:         mediaFileSize.Int64,
				MimeType:         mediaMimeType.String,
				UploadedAt:       mediaUploadedAt.Time,
				URL:              fmt.Sprintf("/api/media/files/%s", mediaFilename.String),
			}
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// GetMessagesByUserID retrieves all messages that a user can access
func (db *DB) GetMessagesByUserID(userID int) ([]models.Message, error) {
	query := `
	SELECT m.id, m.sender_id, m.chat_id, m.content, m.created_at, m.is_read
	FROM messages m
	JOIN chat_members cm ON m.chat_id = cm.chat_id
	WHERE cm.user_id = ?
	ORDER BY m.created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ChatID, &message.Content,
			&message.CreatedAt, &message.IsRead); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	log.Println("Messages retrieved successfully")
	return messages, nil
}

// GetMessagesByUserInChat retrieves all messages sent by a specific user in a chat
func (db *DB) GetMessagesByUserInChat(userID, chatID int) ([]models.Message, error) {
	query := `
	SELECT id, sender_id, chat_id, content, created_at, is_read
	FROM messages
	WHERE sender_id = ? AND chat_id = ?
	ORDER BY created_at DESC`

	rows, err := db.Query(query, userID, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages in chat: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ChatID, &message.Content,
			&message.CreatedAt, &message.IsRead); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	log.Printf("Retrieved %d messages from user %d in chat %d", len(messages), userID, chatID)
	return messages, nil
}

// SearchMessages searches for messages containing specific text
func (db *DB) SearchMessages(chatID int, searchText string) ([]models.Message, error) {
	query := `
	SELECT id, sender_id, chat_id, content, created_at, is_read
	FROM messages
	WHERE chat_id = ? AND content LIKE ?
	ORDER BY created_at DESC`

	// Add wildcards for SQL LIKE
	searchPattern := "%" + searchText + "%"

	rows, err := db.Query(query, chatID, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ChatID, &message.Content,
			&message.CreatedAt, &message.IsRead); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	log.Printf("Found %d messages matching '%s' in chat %d", len(messages), searchText, chatID)
	return messages, nil
}

// MarkMessagesAsRead marks all messages in a chat as read for a user
func (db *DB) MarkMessagesAsRead(userID, chatID int) error {
	query := `
	UPDATE messages
	SET is_read = TRUE
	WHERE chat_id = ? AND sender_id != ? AND is_read = FALSE`

	result, err := db.Exec(query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Update the last_read_at timestamp for the user in this chat
	updateQuery := `
	UPDATE chat_members
	SET last_read_at = ?
	WHERE user_id = ? AND chat_id = ?`

	_, err = db.Exec(updateQuery, time.Now(), userID, chatID)
	if err != nil {
		log.Printf("Warning: failed to update last_read timestamp: %v", err)
	}

	log.Printf("%d messages marked as read", rowsAffected)
	return nil
}

// DeleteMessage deletes a message (if the user is the sender or an admin)
func (db *DB) DeleteMessage(messageID, userID int) error {
	// First check if user is sender or admin
	query := `
	SELECT m.id
	FROM messages m
	LEFT JOIN chat_members cm ON m.chat_id = cm.chat_id AND cm.user_id = ?
	WHERE m.id = ? AND (m.sender_id = ? OR cm.role = 'admin')`

	var id int
	err := db.QueryRow(query, userID, messageID, userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("not authorized to delete this message")
		}
		return fmt.Errorf("failed to check message permissions: %w", err)
	}

	// Delete the message
	deleteQuery := `DELETE FROM messages WHERE id = ?`
	_, err = db.Exec(deleteQuery, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	log.Printf("Message %d deleted by user %d", messageID, userID)
	return nil
}
