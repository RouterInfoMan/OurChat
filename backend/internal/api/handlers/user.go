package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		ID                int        `json:"id"`
		Username          string     `json:"username"`
		Email             string     `json:"email"`
		ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
		Status            string     `json:"status"`
		CreatedAt         time.Time  `json:"created_at"`
		LastLogin         *time.Time `json:"last_login,omitempty"`
	}

	profile := ProfileResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		ProfilePictureURL: user.ProfilePictureURL,
		Status:            user.Status,
		CreatedAt:         user.CreatedAt,
		LastLogin:         user.LastLogin,
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
		ID                int        `json:"id"`
		Username          string     `json:"username"`
		Email             string     `json:"email"`
		ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
		Status            string     `json:"status"`
		CreatedAt         time.Time  `json:"created_at"`
		LastLogin         *time.Time `json:"last_login,omitempty"`
	}

	profile := ProfileResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		ProfilePictureURL: user.ProfilePictureURL,
		Status:            user.Status,
		CreatedAt:         user.CreatedAt,
		LastLogin:         user.LastLogin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// UserIDsRequest represents a request to get user details by IDs
type UserIDsRequest struct {
	UserIDs []int `json:"user_ids"`
}

// HandleGetUsersByIDs gets basic information for a list of user IDs
// Usage: GET /users/ids?ids=1,2,3
// Usage: POST /users/ids with body {"user_ids": [1, 2, 3]}
func (h *UserHandler) HandleGetUsersByIDs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body for POST method
	if r.Method == http.MethodPost {
		var req UserIDsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get users from database
		users, err := h.DB.GetUsersByIDs(req.UserIDs)
		if err != nil {
			log.Printf("Failed to get users: %v", err)
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		return
	}

	// For GET method, parse query parameters
	userIDsParam := r.URL.Query().Get("ids")
	if userIDsParam == "" {
		http.Error(w, "Missing user IDs", http.StatusBadRequest)
		return
	}

	// Split the comma-separated IDs
	idStrings := strings.Split(userIDsParam, ",")
	userIDs := make([]int, 0, len(idStrings))

	for _, idStr := range idStrings {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}
		userIDs = append(userIDs, id)
	}

	// Get users from database
	users, err := h.DB.GetUsersByIDs(userIDs)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// HandleSearchUsers searches for users by partial username match
func (h *UserHandler) HandleSearchUsers(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	currentUserID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get search query from URL parameters
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Validate minimum search length (prevent too broad searches)
	if len(strings.TrimSpace(searchTerm)) < 4 {
		http.Error(w, "Search query must be at least 4 characters long", http.StatusBadRequest)
		return
	}

	// Get optional limit parameter
	limit := 20 // Default limit
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// Search for users, excluding current user
	users, err := h.DB.SearchUsersByName(strings.TrimSpace(searchTerm), limit, currentUserID)
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"count": len(users),
		"query": searchTerm,
	})
}
