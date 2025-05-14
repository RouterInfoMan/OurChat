package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"OurChat/internal/db"
)

// UserHandler contains handlers related to user management
type UserHandler struct {
	DB *db.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *db.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

// HandleGetProfile gets the current user's profile
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	// Create response without sensitive fields
	type ProfileResponse struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Status    string     `json:"status"`
		CreatedAt time.Time  `json:"created_at"`
		LastLogin *time.Time `json:"last_login,omitempty"`
	}

	profile := ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// HandleUpdateProfile updates the current user's profile
func (h *UserHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	type UpdateProfileRequest struct {
		Email  string `json:"email,omitempty"`
		Status string `json:"status,omitempty"`
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate inputs
	updates := make(map[string]interface{})

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Status != "" {
		// Validate status
		validStatuses := map[string]bool{
			"online":  true,
			"offline": true,
			"away":    true,
			"busy":    true,
		}

		if !validStatuses[req.Status] {
			http.Error(w, "Invalid status value", http.StatusBadRequest)
			return
		}

		updates["status"] = req.Status
	}

	// If no updates, return error
	if len(updates) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	// Update profile in database
	if err := h.DB.UpdateUserProfile(userID, updates); err != nil {
		log.Printf("Failed to update user profile: %v", err)

		// Check for specific errors
		if err.Error() == "email is already in use" {
			http.Error(w, "Email address is already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Get updated user profile
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to get updated user profile: %v", err)
		http.Error(w, "Profile updated but failed to retrieve", http.StatusInternalServerError)
		return
	}

	// Create response without sensitive fields
	type ProfileResponse struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Status    string     `json:"status"`
		CreatedAt time.Time  `json:"created_at"`
		LastLogin *time.Time `json:"last_login,omitempty"`
	}

	profile := ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}
