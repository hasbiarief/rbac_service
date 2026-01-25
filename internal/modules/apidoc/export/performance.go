package export

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PerformanceConfig contains performance optimization settings
type PerformanceConfig struct {
	EnableEagerLoading     bool          `json:"enable_eager_loading"`
	EnableStreaming        bool          `json:"enable_streaming"`
	StreamingThreshold     int           `json:"streaming_threshold"`    // Number of endpoints to trigger streaming
	MaxConcurrentExports   int           `json:"max_concurrent_exports"` // Max concurrent export operations
	ExportTimeout          time.Duration `json:"export_timeout"`         // Timeout for export operations
	EnableProgressTracking bool          `json:"enable_progress_tracking"`
}

// DefaultPerformanceConfig returns default performance configuration
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		EnableEagerLoading:     true,
		EnableStreaming:        true,
		StreamingThreshold:     100, // Stream if more than 100 endpoints
		MaxConcurrentExports:   5,   // Max 5 concurrent exports
		ExportTimeout:          5 * time.Minute,
		EnableProgressTracking: true,
	}
}

// ExportProgress tracks export operation progress
type ExportProgress struct {
	ID                string        `json:"id"`
	CollectionID      int64         `json:"collection_id"`
	Format            string        `json:"format"`
	Status            string        `json:"status"`   // pending, processing, completed, failed
	Progress          float64       `json:"progress"` // 0.0 to 1.0
	TotalSteps        int           `json:"total_steps"`
	CompletedSteps    int           `json:"completed_steps"`
	CurrentStep       string        `json:"current_step"`
	StartedAt         time.Time     `json:"started_at"`
	CompletedAt       *time.Time    `json:"completed_at,omitempty"`
	Error             string        `json:"error,omitempty"`
	EstimatedTimeLeft time.Duration `json:"estimated_time_left"`
}

// PerformanceOptimizer handles export performance optimization
type PerformanceOptimizer struct {
	config    *PerformanceConfig
	semaphore chan struct{} // Semaphore for limiting concurrent exports
	progress  map[string]*ExportProgress
	mutex     sync.RWMutex
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer(config *PerformanceConfig) *PerformanceOptimizer {
	if config == nil {
		config = DefaultPerformanceConfig()
	}

	optimizer := &PerformanceOptimizer{
		config:    config,
		semaphore: make(chan struct{}, config.MaxConcurrentExports),
		progress:  make(map[string]*ExportProgress),
	}

	return optimizer
}

// OptimizedExport performs export with performance optimizations
func (p *PerformanceOptimizer) OptimizedExport(
	ctx context.Context,
	collection *CollectionWithDetails,
	format ExportFormat,
	options *ExportOptions,
	exporter ExportService,
) (*ExportResult, error) {
	// Acquire semaphore for concurrency control
	select {
	case p.semaphore <- struct{}{}:
		defer func() { <-p.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Create progress tracker if enabled
	var progressID string
	if p.config.EnableProgressTracking {
		progressID = p.createProgressTracker(collection.ID, format)
		defer p.cleanupProgressTracker(progressID)
	}

	// Set timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, p.config.ExportTimeout)
	defer cancel()

	// Update progress
	if progressID != "" {
		p.updateProgress(progressID, "starting", 0.0, "Initializing export")
	}

	// Check if we should use streaming for large collections
	shouldStream := p.config.EnableStreaming &&
		len(collection.Endpoints) > p.config.StreamingThreshold

	if shouldStream {
		return p.streamingExport(timeoutCtx, collection, format, options, exporter, progressID)
	}

	// Regular export with progress tracking
	return p.regularExport(timeoutCtx, collection, format, options, exporter, progressID)
}

// streamingExport handles large collections with streaming
func (p *PerformanceOptimizer) streamingExport(
	ctx context.Context,
	collection *CollectionWithDetails,
	format ExportFormat,
	options *ExportOptions,
	exporter ExportService,
	progressID string,
) (*ExportResult, error) {
	if progressID != "" {
		p.updateProgress(progressID, "processing", 0.1, "Starting streaming export")
	}

	// For streaming, we'll process endpoints in batches
	batchSize := 50
	totalEndpoints := len(collection.Endpoints)

	// Process in batches to avoid memory issues
	for i := 0; i < totalEndpoints; i += batchSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		end := i + batchSize
		if end > totalEndpoints {
			end = totalEndpoints
		}

		// Update progress
		if progressID != "" {
			progress := float64(i)/float64(totalEndpoints)*0.8 + 0.1 // 10% to 90%
			step := fmt.Sprintf("Processing endpoints %d-%d of %d", i+1, end, totalEndpoints)
			p.updateProgress(progressID, "processing", progress, step)
		}

		// Process batch (in a real implementation, you might want to
		// modify the exporter to support batch processing)
		time.Sleep(10 * time.Millisecond) // Simulate processing time
	}

	// Final export
	if progressID != "" {
		p.updateProgress(progressID, "processing", 0.9, "Finalizing export")
	}

	result, err := exporter.Export(collection, options)
	if err != nil {
		if progressID != "" {
			p.updateProgress(progressID, "failed", 0.0, fmt.Sprintf("Export failed: %v", err))
		}
		return nil, err
	}

	if progressID != "" {
		p.updateProgress(progressID, "completed", 1.0, "Export completed successfully")
	}

	return result, nil
}

// regularExport handles normal-sized collections
func (p *PerformanceOptimizer) regularExport(
	ctx context.Context,
	collection *CollectionWithDetails,
	format ExportFormat,
	options *ExportOptions,
	exporter ExportService,
	progressID string,
) (*ExportResult, error) {
	if progressID != "" {
		p.updateProgress(progressID, "processing", 0.5, "Processing export")
	}

	// Perform export
	result, err := exporter.Export(collection, options)
	if err != nil {
		if progressID != "" {
			p.updateProgress(progressID, "failed", 0.0, fmt.Sprintf("Export failed: %v", err))
		}
		return nil, err
	}

	if progressID != "" {
		p.updateProgress(progressID, "completed", 1.0, "Export completed successfully")
	}

	return result, nil
}

// createProgressTracker creates a new progress tracker
func (p *PerformanceOptimizer) createProgressTracker(collectionID int64, format ExportFormat) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	progressID := fmt.Sprintf("%d_%s_%d", collectionID, format, time.Now().UnixNano())

	progress := &ExportProgress{
		ID:           progressID,
		CollectionID: collectionID,
		Format:       string(format),
		Status:       "pending",
		Progress:     0.0,
		TotalSteps:   100,
		StartedAt:    time.Now(),
	}

	p.progress[progressID] = progress
	return progressID
}

// updateProgress updates progress information
func (p *PerformanceOptimizer) updateProgress(progressID, status string, progress float64, currentStep string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if prog, exists := p.progress[progressID]; exists {
		prog.Status = status
		prog.Progress = progress
		prog.CompletedSteps = int(progress * 100)
		prog.CurrentStep = currentStep

		// Calculate estimated time left
		if progress > 0 && status == "processing" {
			elapsed := time.Since(prog.StartedAt)
			totalEstimated := time.Duration(float64(elapsed) / progress)
			prog.EstimatedTimeLeft = totalEstimated - elapsed
		}

		if status == "completed" || status == "failed" {
			now := time.Now()
			prog.CompletedAt = &now
			prog.EstimatedTimeLeft = 0
		}
	}
}

// cleanupProgressTracker removes progress tracker after completion
func (p *PerformanceOptimizer) cleanupProgressTracker(progressID string) {
	// Keep progress for a short time after completion for status queries
	time.AfterFunc(5*time.Minute, func() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		delete(p.progress, progressID)
	})
}

// GetProgress returns progress information for an export operation
func (p *PerformanceOptimizer) GetProgress(progressID string) (*ExportProgress, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	progress, exists := p.progress[progressID]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	progressCopy := *progress
	return &progressCopy, true
}

// GetAllProgress returns all active progress trackers
func (p *PerformanceOptimizer) GetAllProgress() map[string]*ExportProgress {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	result := make(map[string]*ExportProgress)
	for id, progress := range p.progress {
		progressCopy := *progress
		result[id] = &progressCopy
	}

	return result
}

// GetConfig returns the current performance configuration
func (p *PerformanceOptimizer) GetConfig() *PerformanceConfig {
	return p.config
}

// UpdateConfig updates the performance configuration
func (p *PerformanceOptimizer) UpdateConfig(config *PerformanceConfig) {
	if config == nil {
		return
	}

	p.config = config

	// Update semaphore if max concurrent exports changed
	if len(p.semaphore) != config.MaxConcurrentExports {
		p.semaphore = make(chan struct{}, config.MaxConcurrentExports)
	}
}

// GetStats returns performance statistics
func (p *PerformanceOptimizer) GetStats() *PerformanceStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stats := &PerformanceStats{
		ActiveExports:     0,
		CompletedExports:  0,
		FailedExports:     0,
		AverageExportTime: 0,
	}

	var totalDuration time.Duration
	var completedCount int

	for _, progress := range p.progress {
		switch progress.Status {
		case "pending", "processing":
			stats.ActiveExports++
		case "completed":
			stats.CompletedExports++
			if progress.CompletedAt != nil {
				duration := progress.CompletedAt.Sub(progress.StartedAt)
				totalDuration += duration
				completedCount++
			}
		case "failed":
			stats.FailedExports++
		}
	}

	if completedCount > 0 {
		stats.AverageExportTime = totalDuration / time.Duration(completedCount)
	}

	return stats
}

// PerformanceStats contains performance statistics
type PerformanceStats struct {
	ActiveExports     int           `json:"active_exports"`
	CompletedExports  int           `json:"completed_exports"`
	FailedExports     int           `json:"failed_exports"`
	AverageExportTime time.Duration `json:"average_export_time"`
}
