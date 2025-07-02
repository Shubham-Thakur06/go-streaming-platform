package service

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/repository"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
)

type Service struct {
	config          *config.Config
	mediaRepo       *repository.MediaRepository
	userRepo        *repository.UserRepository
	storageProvider *storage.StorageProvider
}

func New(cfg *config.Config, repo *repository.Repository, storageProvider *storage.StorageProvider) *Service {
	return &Service{
		config:          cfg,
		mediaRepo:       repository.NewMediaRepository(repo),
		userRepo:        repository.NewUserRepository(repo),
		storageProvider: storageProvider,
	}
}
