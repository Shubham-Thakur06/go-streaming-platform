package service

import (
	"io"

	"github.com/google/uuid"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/models"
)

type MediaService struct {
	*Service
}

func NewMediaService(service *Service) *MediaService {
	return &MediaService{Service: service}
}

func (s *MediaService) UploadMedia(file io.Reader, filename, title, description, genre, tags string, userID string, isPublic bool) (*models.Media, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	storageURL, err := s.storageProvider.UploadFile(file, filename, "application/octet-stream")
	if err != nil {
		return nil, err
	}

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

func (s *MediaService) GetMedia(id string) (*models.Media, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return nil, ErrMediaNotFound
	}

	go s.mediaRepo.IncrementViewCount(id)

	return media, nil
}

func (s *MediaService) ListPublicMedia(page, limit int, genre, search string) ([]models.Media, int64, error) {
	return s.mediaRepo.FindPublic(page, limit, genre, search)
}

// DeleteMedia deletes media by ID and user ID
func (s *MediaService) DeleteMedia(id, userID string) error {
	// Delete from storage
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return ErrMediaNotFound
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidCredentials
	}

	if media.UserID != parsedUserID {
		return ErrAccessDenied
	}

	if err := s.storageProvider.DeleteFile(media.Filename); err != nil {
		return err
	}

	return s.mediaRepo.DeleteByIDAndUserID(id, userID)
}

func (s *MediaService) GenerateStreamURL(id string) (string, error) {
	media, err := s.mediaRepo.FindByID(id)
	if err != nil {
		return "", ErrMediaNotFound
	}

	return s.storageProvider.GeneratePresignedURL(media.Filename, 10800) // 3 hour expiry
}
