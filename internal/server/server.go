package server

import (
	"fmt"
	"net/http"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/database"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/handlers"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg             *config.Config
	db              *database.DB
	storageProvider *storage.StorageProvider
	handlers        *handlers.Handlers
	router          *gin.Engine
}

func New(cfg *config.Config, db *database.DB, storageProvider *storage.StorageProvider, handlers *handlers.Handlers) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		cfg:             cfg,
		db:              db,
		storageProvider: storageProvider,
		handlers:        handlers,
		router:          router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	public := s.router.Group("/api/v1")
	{
		public.POST("/auth/login", s.handlers.User.Login)
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.router.Run(addr)
}
