package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"OurChat/internal/db"

	"github.com/gorilla/mux"
)

// ChatHandler contains handlers related to chat management
type ChatHandler struct {
	DB *db.DB
}

// NewChatHandler creates a new chat handler
func NewChatHandler(db *db.DB) *ChatHandler {
	return &ChatHandler{
		DB: db,
	}
}

// HandleGetChats gets all chats for the current user
func (h *ChatHandler) HandleGetChats(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get chats from database
	chats, err := h.DB.GetChatsForUser(userID)
	if err != nil {
		http.Error(w, "Failed to get chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

// ChatRequest represents a request to create a new chat
type ChatRequest struct {
	Type  string `json:"type"` // "direct" or "group"
	Name  string `json:"name,omitempty"`
	Users []int  `json:"users,omitempty"` // For direct chat, should have one user ID
}

// HandleCreateChat creates a new chat
func (h *ChatHandler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Type != "direct" && req.Type != "group" {
		http.Error(w, "Invalid chat type", http.StatusBadRequest)
		return
	}

	if req.Type == "group" && req.Name == "" {
		http.Error(w, "Group chat name is required", http.StatusBadRequest)
		return
	}

	if req.Type == "direct" && (len(req.Users) != 1) {
		http.Error(w, "Direct chat requires exactly one other user", http.StatusBadRequest)
		return
	}

	if len(req.Users) > 0 {
		// Validate each user exists
		for _, otherUserID := range req.Users {
			// Check if user exists
			_, err := h.DB.GetUserByID(userID)
			if err != nil {
				http.Error(w, fmt.Sprintf("User with ID %d does not exist", userID), http.StatusBadRequest)
				return
			}

			// For direct chats, make sure they're not trying to chat with themselves
			if req.Type == "direct" && userID == otherUserID {
				http.Error(w, "Cannot create direct chat with yourself", http.StatusBadRequest)
				return
			}
		}
	}

	// For direct chats, check if a chat already exists
	if req.Type == "direct" {
		otherUserID := req.Users[0]

		// Get or create direct chat
		chatID, err := h.DB.GetDirectChatBetweenUsers(userID, otherUserID)
		if err != nil {
			http.Error(w, "Failed to create direct chat", http.StatusInternalServerError)
			return
		}

		// Get the chat
		chat, err := h.DB.GetChatByID(chatID)
		if err != nil {
			http.Error(w, "Failed to get chat", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
		return
	}

	// For group chats, create a new chat
	chatID, err := h.DB.CreateChat(req.Type, req.Name)
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	// Add current user as admin
	err = h.DB.AddUserToChat(userID, int(chatID), "admin")
	if err != nil {
		http.Error(w, "Failed to add user to chat", http.StatusInternalServerError)
		return
	}

	// Add other users
	for _, otherUserID := range req.Users {
		err = h.DB.AddUserToChat(otherUserID, int(chatID), "member")
		if err != nil {
			// Just log the error and continue
			// TODO: Better error handling
			continue
		}
	}

	// Get the chat
	chat, err := h.DB.GetChatByID(int(chatID))
	if err != nil {
		http.Error(w, "Failed to get chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

// HandleGetChat gets a specific chat
func (h *ChatHandler) HandleGetChat(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get chat ID from URL
	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Check if user is a member of the chat
	isMember, _, err := h.DB.IsUserChatMember(userID, chatID)
	if err != nil {
		http.Error(w, "Failed to verify chat membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "You are not a member of this chat", http.StatusForbidden)
		return
	}

	// Get chat from database
	chat, err := h.DB.GetChatByID(chatID)
	if err != nil {
		http.Error(w, "Failed to get chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// HandleGetChatMembers gets all members of a specific chat
func (h *ChatHandler) HandleGetChatMembers(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get chat ID from URL
	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Check if user is a member of the chat
	isMember, _, err := h.DB.IsUserChatMember(userID, chatID)
	if err != nil {
		http.Error(w, "Failed to verify chat membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "You are not a member of this chat", http.StatusForbidden)
		return
	}

	// Get chat members from database
	members, err := h.DB.GetChatMembers(chatID)
	if err != nil {
		http.Error(w, "Failed to get chat members", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
