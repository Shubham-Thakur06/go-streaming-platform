package storage

import (
	"fmt"
	"io"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/interfaces"
)

type StorageProvider struct {
	provider interfaces.Provider
}

func NewProvider(cfg config.StorageConfig) (*StorageProvider, error) {
	var provider interfaces.Provider
	var err error

	switch cfg.Provider {
	case "aws":
		provider, err = NewAWSProvider(cfg.AWS)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s provider: %w", cfg.Provider, err)
	}

	return &StorageProvider{provider: provider}, nil
}

func (sp *StorageProvider) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	return sp.provider.UploadFile(file, filename, contentType)
}

func (sp *StorageProvider) DeleteFile(filename string) error {
	return sp.provider.DeleteFile(filename)
}

func (sp *StorageProvider) GetFileURL(filename string) (string, error) {
	return sp.provider.GetFileURL(filename)
}

func (sp *StorageProvider) GeneratePresignedURL(filename string, expiresIn int64) (string, error) {
	return sp.provider.GeneratePresignedURL(filename, expiresIn)
}
