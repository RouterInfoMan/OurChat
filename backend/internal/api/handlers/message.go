package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OurChat/internal/db"

	"github.com/gorilla/mux"
)

// MessageHandler contains handlers related to chat messages
type MessageHandler struct {
	DB *db.DB
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(db *db.DB) *MessageHandler {
	return &MessageHandler{
		DB: db,
	}
}

// HandleGetMessages gets messages for a specific chat
func (h *MessageHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
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

	// Get pagination parameters
	limit := 50 // Default limit
	offset := 0 // Default offset

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 50
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}
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

	// Get messages from database
	messages, err := h.DB.GetMessagesByChatID(chatID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// MessageRequest represents a request to send a message
type MessageRequest struct {
	Content string `json:"content"`
}

// HandleSendMessage sends a message to a specific chat
func (h *MessageHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Content == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}

	// TODO: Check if user is a member of the chat

	// Create message
	messageID, err := h.DB.CreateMessage(userID, chatID, req.Content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Get the message
	message, err := h.DB.GetMessageByID(int(messageID))
	if err != nil {
		http.Error(w, "Message sent but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// HandleMarkMessagesAsRead marks all messages in a chat as read
func (h *MessageHandler) HandleMarkMessagesAsRead(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Check if user is a member of the chat

	// Mark messages as read
	err = h.DB.MarkMessagesAsRead(userID, chatID)
	if err != nil {
		http.Error(w, "Failed to mark messages as read", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Messages marked as read",
	})
}

// HandleSearchMessages searches for messages in a chat
func (h *MessageHandler) HandleSearchMessages(w http.ResponseWriter, r *http.Request) {
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

	// Get search query
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
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

	// Search messages
	messages, err := h.DB.SearchMessages(chatID, query)
	if err != nil {
		http.Error(w, "Failed to search messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
