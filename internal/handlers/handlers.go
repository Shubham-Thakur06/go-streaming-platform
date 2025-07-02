package handlers

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
	"gorm.io/gorm"
)

type Handlers struct {
	User  *UserHandler
	Media *MediaHandler
}

type MediaHandler struct {
	db              *gorm.DB
	cfg             *config.Config
	storageProvider *storage.StorageProvider
}

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}
