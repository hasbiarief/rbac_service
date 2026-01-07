package jobs

import (
	"gin-scalable-api/internal/service"
	"log"
	"time"
)

type CleanupJob struct {
	authService *service.AuthService
	interval    time.Duration
	stopChan    chan bool
}

func NewCleanupJob(authService *service.AuthService, interval time.Duration) *CleanupJob {
	return &CleanupJob{
		authService: authService,
		interval:    interval,
		stopChan:    make(chan bool),
	}
}

// Start begins the cleanup job
func (c *CleanupJob) Start() {
	log.Println("Starting token cleanup job...")

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// Run cleanup immediately on start
	c.runCleanup()

	for {
		select {
		case <-ticker.C:
			c.runCleanup()
		case <-c.stopChan:
			log.Println("Token cleanup job stopped")
			return
		}
	}
}

// Stop stops the cleanup job
func (c *CleanupJob) Stop() {
	c.stopChan <- true
}

func (c *CleanupJob) runCleanup() {
	log.Println("Running token cleanup...")

	if err := c.authService.CleanupExpiredTokens(); err != nil {
		log.Printf("Error during token cleanup: %v", err)
	} else {
		log.Println("Token cleanup completed successfully")
	}
}
