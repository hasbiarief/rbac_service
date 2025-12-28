package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gin-scalable-api/config"
	"gin-scalable-api/pkg/database"
	"gin-scalable-api/pkg/migration"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, status")
		dir    = flag.String("dir", "migrations", "Migrations directory")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator := migration.NewMigrator(db.DB, *dir)

	// Execute action
	switch *action {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, status")
		os.Exit(1)
	}
}
