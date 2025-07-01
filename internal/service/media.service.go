package service

import (
	"io"

	"github.com/google/uuid"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/repository"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
)

// Service provides business logic methods
type Service struct {
	config          *config.Config
	mediaRepo       *repository.MediaRepository
	storageProvider *storage.StorageProvider
}

// New creates a new service instance
func New(cfg *config.Config, repo *repository.Repository, storageProvider *storage.StorageProvider) *Service {
	return &Service{
		config:          cfg,
		mediaRepo:       repository.NewMediaRepository(repo),
		storageProvider: storageProvider,
	}
}

// MediaService provides media-related business logic
type MediaService struct {
	*Service
}

// NewMediaService creates a new media service
func NewMediaService(service *Service) *MediaService {
	return &MediaService{Service: service}
}

// UploadMedia uploads a media file and creates a media record
func (s *MediaService) UploadMedia(file io.Reader, filename, title, description, genre, tags string, userID string, isPublic bool) (*models.Media, error) {
	// Parse user ID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Upload file to storage
	storageURL, err := s.storageProvider.UploadFile(file, filename, "application/octet-stream")
	if err != nil {
		return nil, err
	}

	// Create media record
	media := &models.Media{
		Title:       title,
		Description: description,
		Filename:    filename,
		Genre:       genre,
		Tags:        tags,
		StorageURL:  storageURL,
		UserID:      parsedUserID,
		IsPublic:    isPublic,
	}

	if err := s.mediaRepo.Create(media); err != nil {
		return nil, err
	}

	return media, nil
}

// GetMedia retrieves media by ID
func (s *MediaService) GetMedia(id string) (*models.Media, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Increment view count
	go s.mediaRepo.IncrementViewCount(id)

	return media, nil
}

// ListPublicMedia retrieves public media with filters
func (s *MediaService) ListPublicMedia(page, limit int, genre, search string) ([]models.Media, int64, error) {
	return s.mediaRepo.FindPublic(page, limit, genre, search)
}

// DeleteMedia deletes media by ID and user ID
func (s *MediaService) DeleteMedia(id, userID string) error {
	// Delete from storage
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.storageProvider.DeleteFile(media.Filename); err != nil {
		return err
	}

	// Delete from database
	return s.mediaRepo.DeleteByIDAndUserID(id, userID)
}

// GenerateStreamURL generates a presigned URL for streaming
func (s *MediaService) GenerateStreamURL(id string) (string, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return "", err
	}

	return s.storageProvider.GeneratePresignedURL(media.Filename, 10800) // 3 hour expiry
}
