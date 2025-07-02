package app

import (
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/database"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/handlers"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/server"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"
)

type App struct {
	Config          *config.Config
	DB              *database.DB
	StorageProvider *storage.StorageProvider
	Server          *server.Server
	Handlers        *handlers.Handlers
}

func New(cfg *config.Config) (*App, error) {
	db, err := database.InitWithHost(cfg.Database, cfg.Host)
	if err != nil {
		return nil, err
	}

	storageProvider, err := storage.NewProvider(cfg.Storage)
	if err != nil {
		return nil, err
	}

	handlers := handlers.New(db, cfg, storageProvider)

	srv := server.New(cfg, db, storageProvider, handlers)

	app := &App{
		Config:          cfg,
		DB:              db,
		StorageProvider: storageProvider,
		Server:          srv,
		Handlers:        handlers,
	}

	return app, nil
}

func (a *App) Run() error {
	return a.Server.Start()
}

func (a *App) Shutdown() error {
	return nil
}
