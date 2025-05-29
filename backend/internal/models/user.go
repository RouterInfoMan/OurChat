package models

import (
	"time"
)

type User struct {
	ID                int        `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	Password          string     `json:"-"` // Never include password in JSON responses
	JWTKey            string     `json:"-"` // Never include JWT key in JSON responses
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	Status            string     `json:"status"`
}

// UserBasic represents basic user information for public endpoints
type UserBasic struct {
	ID                int     `json:"id"`
	Username          string  `json:"username"`
	Status            string  `json:"status"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
}
type UserProfile struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Status    string     `json:"status"`
}
