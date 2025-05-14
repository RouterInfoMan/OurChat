package db

import (
	"database/sql"
	"fmt"
)

// IsUserChatMember checks if a user is a member of a specific chat
func (db *DB) IsUserChatMember(userID, chatID int) (bool, string, error) {
	query := `
	SELECT role FROM chat_members
	WHERE user_id = ? AND chat_id = ?`

	var role string
	err := db.QueryRow(query, userID, chatID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			// User is not a member of the chat
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check chat membership: %w", err)
	}

	// User is a member of the chat
	return true, role, nil
}
