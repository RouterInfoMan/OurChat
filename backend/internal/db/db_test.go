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

func removeFile(dbPath string) error {
	return os.Remove(dbPath)
}
