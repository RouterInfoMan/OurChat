package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"OurChat/internal/db"
	"OurChat/internal/models"

	"github.com/gorilla/mux"
)

// MediaHandler handles media file uploads and serving
type MediaHandler struct {
	DB        *db.DB
	UploadDir string
}

// NewMediaHandler creates a new media handler
func NewMediaHandler(db *db.DB) *MediaHandler {
	uploadDir := "./uploads"
	// Create upload directory if it doesn't exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(filepath.Join(uploadDir, "profiles"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "media"), 0755)

	return &MediaHandler{
		DB:        db,
		UploadDir: uploadDir,
	}
}

// HandleUploadProfilePicture handles profile picture uploads
func (h *MediaHandler) HandleUploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (5MB max)
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Header.Get("Content-Type")) {
		http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF are allowed", http.StatusBadRequest)
		return
	}

	// Validate file size (5MB max)
	if header.Size > 5*1024*1024 {
		http.Error(w, "File too large. Maximum size is 5MB", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename, userID)
	filePath := filepath.Join(h.UploadDir, "profiles", filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Update user profile picture in database
	profileURL := fmt.Sprintf("/api/media/profiles/%s", filename)
	err = h.DB.UpdateUserProfilePicture(userID, profileURL)
	if err != nil {
		// Clean up file if database update fails
		os.Remove(filePath)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile picture updated successfully",
		"url":     profileURL,
	})
}

// HandleUploadMedia handles general media file uploads
func (h *MediaHandler) HandleUploadMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (50MB max)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidMediaType(header.Header.Get("Content-Type")) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Validate file size (50MB max)
	if header.Size > 50*1024*1024 {
		http.Error(w, "File too large. Maximum size is 50MB", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename, userID)
	filePath := filepath.Join(h.UploadDir, "media", filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
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
		// Clean up file if database save fails
		os.Remove(filePath)
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}

	mediaFile.ID = int(mediaFileID)
	mediaFile.URL = fmt.Sprintf("/api/media/files/%s", filename)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mediaFile)
}

// HandleServeMedia serves uploaded media files with authorization
func (h *MediaHandler) HandleServeMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]
	mediaType := vars["type"]

	if filename == "" || mediaType == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var filePath string
	var authorized bool
	var err error

	switch mediaType {
	case "profiles":
		authorized, filePath, err = h.checkProfilePictureAccess(userID, filename)
	case "files":
		authorized, filePath, err = h.checkMediaFileAccess(userID, filename)
	default:
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to verify access", http.StatusInternalServerError)
		return
	}

	if !authorized {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// Helper functions
func (h *MediaHandler) checkProfilePictureAccess(userID int, filename string) (bool, string, error) {
	filePath := filepath.Join(h.UploadDir, "profiles", filename)

	// Get the owner of this profile picture
	ownerID, err := h.DB.GetProfilePictureOwner(filename)
	if err != nil {
		return false, "", err
	}

	// User can access their own profile picture
	if ownerID == userID {
		return true, filePath, nil
	}

	// User can access profile pictures of people they share chats with
	hasSharedChat, err := h.DB.UsersShareChat(userID, ownerID)
	if err != nil {
		return false, "", err
	}

	return hasSharedChat, filePath, nil
}

func (h *MediaHandler) checkMediaFileAccess(userID int, filename string) (bool, string, error) {
	filePath := filepath.Join(h.UploadDir, "media", filename)

	// Get media file info by filename
	mediaFile, err := h.DB.GetMediaFileByFilename(filename)
	if err != nil {
		return false, "", err
	}

	// User can access their own uploaded files
	if mediaFile.UploadedBy == userID {
		return true, filePath, nil
	}

	// Check if the media file was shared in a chat that the user is a member of
	hasAccess, err := h.DB.UserHasAccessToMediaFile(userID, mediaFile.ID)
	if err != nil {
		return false, "", err
	}

	return hasAccess, filePath, nil
}

func isValidImageType(mimeType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
	}

	for _, validType := range validTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}

func isValidMediaType(mimeType string) bool {
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

func generateUniqueFilename(originalFilename string, userID int) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)

	// Create hash from user ID, timestamp, and original filename
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%d_%s", userID, time.Now().UnixNano(), originalFilename)))

	// Generate filename: hash + extension
	return fmt.Sprintf("%x%s", hash, ext)
}
