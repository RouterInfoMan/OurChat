package models

import "time"

// Message represents a chat message
type Message struct {
	ID          int        `json:"id"`
	SenderID    int        `json:"sender_id"`
	ChatID      int        `json:"chat_id"`
	Content     string     `json:"content"`
	MessageType string     `json:"message_type"` // "text" or "media"
	MediaFileID *int       `json:"media_file_id,omitempty"`
	MediaFile   *MediaFile `json:"media_file,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	IsRead      bool       `json:"is_read"`
}
