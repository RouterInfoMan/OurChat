package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	// Temporary database file for testing
	dbPath := "./test.db"
	defer func() {
		// Clean up test database after the test is done
		if err := removeFile(dbPath); err != nil {
			t.Errorf("failed to delete test db: %v", err)
		}
	}()

	// Create a new DB instance
	database, err := NewDB(dbPath)
	assert.NoError(t, err)
	assert.NotNil(t, database)

	// Test if the table was created successfully
	_, err = database.GetUserByUsername("nonexistent") // Should not return an error, just empty result
	assert.Error(t, err)
}

func TestCreateUser(t *testing.T) {
	// Temporary database file for testing
	dbPath := "./test.db"
	defer func() {
		// Clean up test database after the test is done
		if err := removeFile(dbPath); err != nil {
			t.Errorf("failed to delete test db: %v", err)
		}
	}()

	// Create a new DB instance
	database, err := NewDB(dbPath)
	assert.NoError(t, err)

	// Create a test user
	err = database.CreateUser("testuser", "testuser@example.com", "password123")
	assert.NoError(t, err)

	// Fetch user by username
	user, err := database.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "testuser@example.com", user.Email)

	// Try to insert the same user again (should fail)
	err = database.CreateUser("testuser", "testuser@example.com", "password123")
	assert.Error(t, err, "expected error when creating duplicate user")
}

func TestGetUserByEmail(t *testing.T) {
	// Temporary database file for testing
	dbPath := "./test.db"
	defer func() {
		// Clean up test database after the test is done
		if err := removeFile(dbPath); err != nil {
			t.Errorf("failed to delete test db: %v", err)
		}
	}()

	// Create a new DB instance
	database, err := NewDB(dbPath)
	assert.NoError(t, err)

	// Create a test user
	err = database.CreateUser("testuser", "testuser@example.com", "password123")
	assert.NoError(t, err)

	// Fetch user by email
	user, err := database.GetUserByEmail("testuser@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "testuser@example.com", user.Email)
}

func TestGetMessagesByUserID(t *testing.T) {
	// Temporary database file for testing
	dbPath := "./test.db"
	defer func() {
		// Clean up test database after the test is done
		if err := removeFile(dbPath); err != nil {
			t.Errorf("failed to delete test db: %v", err)
		}
	}()

	// Create a new DB instance
	database, err := NewDB(dbPath)
	assert.NoError(t, err)

	// Create a test user
	err = database.CreateUser("testuser", "testuser@example.com", "password123")
	assert.NoError(t, err)

	// Create a test message
	err = database.CreateMessage(1, 1, "hello, world!")
	assert.NoError(t, err)

	// Fetch messages by user ID
	messages, err := database.GetMessagesByUserID(1)
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "hello, world!", messages[0].Content)
}
func TestGetMessagesByUserToUser(t *testing.T) {
	// Temporary database file for testing
	dbPath := "./test.db"
	defer func() {
		// Clean up test database after the test is done
		if err := removeFile(dbPath); err != nil {
			t.Errorf("failed to delete test db: %v", err)
		}
	}()

	// Create a new DB instance
	database, err := NewDB(dbPath)
	assert.NoError(t, err)

	// Create test users
	err = database.CreateUser("user1", "user1@example.com", "password123")
	assert.NoError(t, err)
	err = database.CreateUser("user2", "user2@example.com", "password123")
	assert.NoError(t, err)

	// Create test messages between the users
	err = database.CreateMessage(1, 2, "Hello from user1 to user2!")
	assert.NoError(t, err)
	err = database.CreateMessage(2, 1, "Hello from user2 to user1!")
	assert.NoError(t, err)

	// Fetch messages between user1 and user2
	messages, err := database.GetMessagesByUserToUser(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Equal(t, 2, len(messages))

	// Validate the content of the messages
	assert.Equal(t, "Hello from user1 to user2!", messages[0].Content)
	assert.Equal(t, "Hello from user2 to user1!", messages[1].Content)
}
func removeFile(dbPath string) error {
	return os.Remove(dbPath)
}
