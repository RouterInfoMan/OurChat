package models

import (
	"time"
)

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	JWTKey    string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
	Status    string     `json:"status"`
}

type UserProfile struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Status    string     `json:"status"`
}

type UserBasic struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}
