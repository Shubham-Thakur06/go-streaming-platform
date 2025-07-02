package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/middleware"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/service"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
	"github.com/gin-gonic/gin"
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

func (h *MediaHandler) UploadMedia(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUnauthorized.Error()})
		return
	}

	var req UploadMediaRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validate file type
	allowedTypes := []string{".mp3", ".mp4", ".wav", ".avi", ".mov", ".mkv"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	filename := fmt.Sprintf("%s/%s", user.ID.String(), file.Filename)

	src, err := file.Open()
	if err != nil {
		fmt.Println("Failed to open file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	storageURL, err := h.storageProvider.UploadFile(src, filename, contentType)
	if err != nil {
		fmt.Println("Failed to upload file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Create media record
	media := models.Media{
		Title:       req.Title,
		Description: req.Description,
		Filename:    filename,
		FileSize:    file.Size,
		FileType:    strings.TrimPrefix(ext, "."),
		Genre:       req.Genre,
		Tags:        req.Tags,
		StorageURL:  storageURL,
		UserID:      user.ID,
		IsPublic:    req.IsPublic,
	}

	if err := h.db.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Media uploaded successfully",
		"media":   h.toMediaResponse(&media),
	})
}

func (h *MediaHandler) GetMedia(c *gin.Context) {
	mediaID := c.Param("id")
	id, err := uuid.Parse(mediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var media models.Media
	if err := h.db.Preload("User").Where("id = ?", id).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrMediaNotFound.Error()})
		return
	}

	// Increment view count
	h.db.Model(&media).Update("view_count", media.ViewCount+1)

	c.JSON(http.StatusOK, gin.H{
		"media": h.toMediaResponse(&media),
	})
}

func (h *MediaHandler) ListMedia(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	genre := c.Query("genre")
	search := c.Query("search")

	offset := (page - 1) * limit
	query := h.db.Model(&models.Media{}).Preload("User").Where("is_public = ?", true)

	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var media []models.Media
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}

	var response []MediaResponse
	for _, m := range media {
		response = append(response, h.toMediaResponse(&m))
	}

	c.JSON(http.StatusOK, gin.H{
		"media": response,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *MediaHandler) GetStreamURL(c *gin.Context) {
	mediaID := c.Param("id")
	id, err := uuid.Parse(mediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var media models.Media
	if err := h.db.Where("id = ?", id).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrMediaNotFound.Error()})
		return
	}

	// Generate presigned URL for streaming (1 hour expiry)
	streamURL, err := h.storageProvider.GeneratePresignedURL(media.Filename, 3600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate stream URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stream_url": streamURL,
		"media":      h.toMediaResponse(&media),
	})
}

func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUnauthorized.Error()})
		return
	}

	mediaID := c.Param("id")
	id, err := uuid.Parse(mediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var media models.Media
	if err := h.db.Where("id = ? AND user_id = ?", id, user.ID).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrMediaNotFound.Error()})
		return
	}

	// Delete from cloud storage
	if err := h.storageProvider.DeleteFile(media.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	// Delete from database
	if err := h.db.Delete(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Media deleted successfully",
	})
}

func (h *MediaHandler) toMediaResponse(media *models.Media) MediaResponse {
	return MediaResponse{
		ID:           media.ID,
		Title:        media.Title,
		Description:  media.Description,
		Filename:     media.Filename,
		FileSize:     media.FileSize,
		FileType:     media.FileType,
		Genre:        media.Genre,
		Tags:         media.Tags,
		Duration:     media.Duration,
		StorageURL:   media.StorageURL,
		ThumbnailURL: media.ThumbnailURL,
		IsPublic:     media.IsPublic,
		ViewCount:    media.ViewCount,
		UserID:       media.UserID,
		CreatedAt:    media.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
