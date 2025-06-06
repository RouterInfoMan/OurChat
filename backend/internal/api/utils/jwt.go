package utils

import (
	"errors"
	"fmt"
	"time"

	"OurChat/internal/db"
	"OurChat/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWT expiration time (24 hours)
const JWTExpiration = time.Hour * 24

// GenerateJWT generates a JWT token for a user
func GenerateJWT(user *models.User) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(JWTExpiration).Unix(),
		"iat":     time.Now().Unix(),
	})

	// Sign token with user's JWT key
	tokenString, err := token.SignedString([]byte(user.JWTKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims if valid
func ValidateJWT(tokenString string, db *db.DB) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		// Extract user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid user_id in token")
		}
		userID := int(userIDFloat)

		// Get the user from the database to retrieve their JWT key
		user, err := db.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		// Return the user's JWT key for validation
		return []byte(user.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// Password reset token expiration (30 minutes)
const PasswordResetExpiration = time.Minute * 30

// GeneratePasswordResetToken creates a JWT token for password reset
func GeneratePasswordResetToken(user *models.User) (string, error) {
	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID,
		"purpose": "password_reset", // Explicit purpose for additional security
		"exp":     time.Now().Add(PasswordResetExpiration).Unix(),
		"iat":     time.Now().Unix(),
	})

	// Sign token with user's current JWT key
	tokenString, err := token.SignedString([]byte(user.JWTKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidatePasswordResetToken validates a password reset token
func ValidatePasswordResetToken(tokenString string, db *db.DB) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		// Check purpose claim
		purpose, ok := claims["purpose"].(string)
		if !ok || purpose != "password_reset" {
			return nil, fmt.Errorf("invalid token purpose")
		}

		// Extract user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid user_id in token")
		}
		userID := int(userIDFloat)

		// Get the user from the database to retrieve their JWT key
		user, err := db.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		// Return the user's JWT key for validation
		return []byte(user.JWTKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Get user ID from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}

	return int(userIDFloat), nil
}
