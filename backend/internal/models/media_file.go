package models

import "time"

// MediaFile represents a media file in the database
type MediaFile struct {
	ID               int       `json:"id"`
	Filename         string    `json:"filename"`
	OriginalFilename string    `json:"original_filename"`
	FilePath         string    `json:"-"` // Don't expose server path
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	UploadedBy       int       `json:"uploaded_by"`
	UploadedAt       time.Time `json:"uploaded_at"`
	URL              string    `json:"url"` // Generated when serving
}
