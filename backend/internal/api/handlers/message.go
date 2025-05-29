package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"OurChat/internal/db"
	"OurChat/internal/models"

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

func (h *MessageHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	limit := 50
	offset := 0

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

	// Get messages with media file information
	messages, err := h.DB.GetMessagesByChatIDWithMedia(chatID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// HandleSendMessage sends a message to a specific chat
type MessageRequest struct {
	Content     string `json:"content,omitempty"`
	MessageType string `json:"message_type"`
	MediaFileID *int   `json:"media_file_id,omitempty"`
}

// REPLACE your existing HandleSendMessage function with this:
func (h *MessageHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set default message type
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	// Validate based on message type
	if req.MessageType == "text" {
		if strings.TrimSpace(req.Content) == "" {
			http.Error(w, "Text message content is required", http.StatusBadRequest)
			return
		}
	} else if req.MessageType == "media" {
		if req.MediaFileID == nil {
			http.Error(w, "Media file ID is required for media messages", http.StatusBadRequest)
			return
		}

		// Verify media file exists and belongs to user
		mediaFile, err := h.DB.GetMediaFileByID(*req.MediaFileID)
		if err != nil {
			http.Error(w, "Media file not found", http.StatusBadRequest)
			return
		}

		if mediaFile.UploadedBy != userID {
			http.Error(w, "You can only send media files you uploaded", http.StatusForbidden)
			return
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

	// Create message
	messageID, err := h.DB.CreateMessage(userID, chatID, req.Content, req.MessageType, req.MediaFileID)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Get the message with media file info
	message, err := h.DB.GetMessageByIDWithMedia(int(messageID))
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
func (h *MessageHandler) HandleSendMediaMessage(w http.ResponseWriter, r *http.Request) {
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

	// Parse multipart form (50MB max)
	err = r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get optional caption
	caption := r.FormValue("caption")

	// Get the file
	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "No media file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !h.isValidMediaType(header.Header.Get("Content-Type")) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Validate file size (50MB max)
	if header.Size > 50*1024*1024 {
		http.Error(w, "File too large. Maximum size is 50MB", http.StatusBadRequest)
		return
	}

	// Save the file and create media record
	mediaFileID, err := h.saveMediaFile(file, header, userID)
	if err != nil {
		http.Error(w, "Failed to save media file", http.StatusInternalServerError)
		return
	}

	// Create the message with media
	messageID, err := h.DB.CreateMessage(userID, chatID, caption, "media", &mediaFileID)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Get the complete message
	message, err := h.DB.GetMessageByIDWithMedia(int(messageID))
	if err != nil {
		http.Error(w, "Message sent but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// Helper functions for the message handler
func (h *MessageHandler) saveMediaFile(file multipart.File, header *multipart.FileHeader, userID int) (int, error) {
	uploadDir := "./uploads"
	os.MkdirAll(filepath.Join(uploadDir, "media"), 0755)

	// Generate unique filename
	filename := h.generateUniqueFilename(header.Filename, userID)
	filePath := filepath.Join(uploadDir, "media", filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	// Reset file pointer to beginning
	file.Seek(0, 0)

	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return 0, err
	}

	// Save file metadata to database
	mediaFile := &models.MediaFile{
		Filename:         filename,
		OriginalFilename: header.Filename,
		FilePath:         filePath,
		FileSize:         header.Size,
		MimeType:         header.Header.Get("Content-Type"),
		UploadedBy:       userID,
		UploadedAt:       time.Now(),
	}

	mediaFileID, err := h.DB.CreateMediaFile(mediaFile)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return 0, err
	}

	return int(mediaFileID), nil
}

func (h *MessageHandler) isValidMediaType(mimeType string) bool {
	validTypes := []string{
		// Images
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
		// Videos
		"video/mp4", "video/webm", "video/avi", "video/mov",
		// Audio
		"audio/mp3", "audio/wav", "audio/ogg", "audio/mpeg",
		// Documents
		"application/pdf", "text/plain",
	}

	for _, validType := range validTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}

func (h *MessageHandler) generateUniqueFilename(originalFilename string, userID int) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)

	// Create hash from user ID, timestamp, and original filename
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%d_%s", userID, time.Now().UnixNano(), originalFilename)))

	// Generate filename: hash + extension
	return fmt.Sprintf("%x%s", hash, ext)
}
