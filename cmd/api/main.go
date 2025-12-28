package main

import (
	"gin-scalable-api/config"
	"gin-scalable-api/internal/server"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize server
	srv := server.NewServer(cfg)

	// Initialize all components (database, repositories, services, handlers, routes)
	if err := srv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
