package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
	*sql.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=60",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection with retry
	var pingErr error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", i+1, maxRetries)
		pingErr = db.Ping()
		if pingErr == nil {
			log.Println("✅ Database connection established successfully")
			return &DB{db}, nil
		}
		log.Printf("❌ Ping attempt %d failed: %v", i+1, pingErr)
		if i < maxRetries-1 {
			time.Sleep(3 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, pingErr)
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Transaction helper
func (db *DB) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
