package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/app"
	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	cfg := config.New()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	if err := application.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped")
}
