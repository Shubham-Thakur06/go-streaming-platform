package handlers

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/database"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
)

type Handlers struct {
	User  *UserHandler
	Media *MediaHandler
}

func New(db *database.DB, cfg *config.Config, storageProvider *storage.StorageProvider) *Handlers {
	return &Handlers{
		User:  NewUserHandler(db.DB, cfg),
		Media: NewMediaHandler(db.DB, cfg, storageProvider),
	}
}
