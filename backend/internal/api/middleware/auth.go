package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"OurChat/internal/api/utils"
	"OurChat/internal/db"
)

// AuthMiddleware is a middleware for JWT authentication
type AuthMiddleware struct {
	DB *db.DB
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(db *db.DB) *AuthMiddleware {
	return &AuthMiddleware{
		DB: db,
	}
}

// Middleware is the JWT authentication middleware handler
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in format: Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims, err := utils.ValidateJWT(tokenString, m.DB)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Get user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
			return
		}

		// Create a new context with the user ID
		ctx := context.WithValue(r.Context(), "user_id", int(userID))

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware wrapper for routes that require authentication
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Middleware(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}
