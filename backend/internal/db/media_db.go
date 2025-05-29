package db

import (
	"database/sql"
	"fmt"

	"OurChat/internal/models"
)

// CreateMediaFile saves media file metadata to the database
func (db *DB) CreateMediaFile(mediaFile *models.MediaFile) (int64, error) {
	query := `
	INSERT INTO media_files (filename, original_filename, file_path, file_size, mime_type, uploaded_by, uploaded_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query,
		mediaFile.Filename,
		mediaFile.OriginalFilename,
		mediaFile.FilePath,
		mediaFile.FileSize,
		mediaFile.MimeType,
		mediaFile.UploadedBy,
		mediaFile.UploadedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create media file: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// GetMediaFileByID retrieves a media file by its ID
func (db *DB) GetMediaFileByID(mediaFileID int) (*models.MediaFile, error) {
	mediaFile := &models.MediaFile{}
	query := `
	SELECT id, filename, original_filename, file_path, file_size, mime_type, uploaded_by, uploaded_at
	FROM media_files WHERE id = ?`

	err := db.QueryRow(query, mediaFileID).Scan(
		&mediaFile.ID,
		&mediaFile.Filename,
		&mediaFile.OriginalFilename,
		&mediaFile.FilePath,
		&mediaFile.FileSize,
		&mediaFile.MimeType,
		&mediaFile.UploadedBy,
		&mediaFile.UploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("media file not found")
		}
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}

	return mediaFile, nil
}

// GetMediaFileByFilename retrieves a media file by its filename
func (db *DB) GetMediaFileByFilename(filename string) (*models.MediaFile, error) {
	mediaFile := &models.MediaFile{}
	query := `
	SELECT id, filename, original_filename, file_path, file_size, mime_type, uploaded_by, uploaded_at
	FROM media_files WHERE filename = ?`

	err := db.QueryRow(query, filename).Scan(
		&mediaFile.ID,
		&mediaFile.Filename,
		&mediaFile.OriginalFilename,
		&mediaFile.FilePath,
		&mediaFile.FileSize,
		&mediaFile.MimeType,
		&mediaFile.UploadedBy,
		&mediaFile.UploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("media file not found")
		}
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}

	return mediaFile, nil
}

// UserHasAccessToMediaFile checks if a user has access to a media file
func (db *DB) UserHasAccessToMediaFile(userID, mediaFileID int) (bool, error) {
	query := `
	SELECT COUNT(*) FROM messages m
	JOIN chat_members cm ON m.chat_id = cm.chat_id
	WHERE m.media_file_id = ? AND cm.user_id = ?`

	var count int
	err := db.QueryRow(query, mediaFileID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check media file access: %w", err)
	}

	return count > 0, nil
}

// GetProfilePictureOwner gets the user ID who owns a profile picture by filename
func (db *DB) GetProfilePictureOwner(filename string) (int, error) {
	profileURL := fmt.Sprintf("/api/media/profiles/%s", filename)

	var userID int
	query := `SELECT id FROM users WHERE profile_picture_url = ?`

	err := db.QueryRow(query, profileURL).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("profile picture owner not found")
		}
		return 0, fmt.Errorf("failed to get profile picture owner: %w", err)
	}

	return userID, nil
}

// UsersShareChat checks if two users share any chat
func (db *DB) UsersShareChat(userID1, userID2 int) (bool, error) {
	query := `
	SELECT COUNT(*) FROM chat_members cm1
	JOIN chat_members cm2 ON cm1.chat_id = cm2.chat_id
	WHERE cm1.user_id = ? AND cm2.user_id = ? AND cm1.user_id != cm2.user_id`

	var count int
	err := db.QueryRow(query, userID1, userID2).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check shared chats: %w", err)
	}

	return count > 0, nil
}

// UpdateUserProfilePicture updates a user's profile picture URL
func (db *DB) UpdateUserProfilePicture(userID int, profilePictureURL string) error {
	query := `UPDATE users SET profile_picture_url = ? WHERE id = ?`
	_, err := db.Exec(query, profilePictureURL, userID)
	if err != nil {
		return fmt.Errorf("failed to update profile picture: %w", err)
	}
	return nil
}
