package models

import (
	"time"
)

type ChatMember struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	ChatID     int        `json:"chat_id"`
	Role       string     `json:"role"`
	JoinedAt   time.Time  `json:"joined_at"`
	LastReadAt *time.Time `json:"last_read_at,omitempty"`

	// Additional fields from the user table
	Username          string  `json:"username"`
	Status            string  `json:"status"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
}
