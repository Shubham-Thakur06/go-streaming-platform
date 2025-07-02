package handlers

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaHandler struct {
	db              *gorm.DB
	cfg             *config.Config
	storageProvider *storage.StorageProvider
}

func NewMediaHandler(db *gorm.DB, cfg *config.Config, storageProvider *storage.StorageProvider) *MediaHandler {
	return &MediaHandler{
		db:              db,
		cfg:             cfg,
		storageProvider: storageProvider,
	}
}

type UploadMediaRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Genre       string `form:"genre"`
	Tags        string `form:"tags"`
	IsPublic    bool   `form:"is_public"`
}

type MediaResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Filename     string    `json:"filename"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	Genre        string    `json:"genre"`
	Tags         string    `json:"tags"`
	Duration     int       `json:"duration"`
	StorageURL   string    `json:"storage_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	IsPublic     bool      `json:"is_public"`
	ViewCount    int       `json:"view_count"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    string    `json:"created_at"`
}
