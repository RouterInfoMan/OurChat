package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"OurChat/internal/api/utils"
	"OurChat/internal/db"

	"golang.org/x/crypto/bcrypt"
)

// AuthHandler contains handlers related to authentication
type AuthHandler struct {
	DB *db.DB
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(db *db.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token   string `json:"token"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// HandleLogin handles user login and token generation
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Login failed for user %s: %v", req.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Verify password (assuming password is hashed with bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Invalid password for user %s", req.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Update last login time
	if err := h.DB.UpdateLastLogin(user.ID); err != nil {
		// Non-critical error, just log it
		log.Printf("Failed to update last login: %v", err)
	}

	// Return token in response
	response := LoginResponse{
		Token:   token,
		UserID:  user.ID,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update user status to offline
	if err := h.DB.UpdateUserStatus(userID, "offline"); err != nil {
		log.Printf("Failed to update user status: %v", err)
		// Continue despite error - not critical
	}

	// Optionally, update JWT key to invalidate all tokens
	// This would force the user to login again on all devices
	// Uncomment if you want this behavior
	/*
		if _, err := h.DB.UpdateJWTKey(userID); err != nil {
			log.Printf("Failed to update JWT key: %v", err)
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// HandleRegister handles user registration
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	_, err := h.DB.GetUserByUsername(req.Username)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Check if email already exists
	_, err = h.DB.GetUserByEmail(req.Email)
	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Create user
	if err := h.DB.CreateUser(req.Username, req.Email, string(hashedPassword)); err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get the created user to get the ID
	user, err := h.DB.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Failed to get created user: %v", err)
		http.Error(w, "User created but failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := RegisterResponse{
		UserID:  user.ID,
		Message: "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
