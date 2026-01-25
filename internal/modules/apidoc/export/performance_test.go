package export

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPerformanceConfig(t *testing.T) {
	config := DefaultPerformanceConfig()

	// Test default values
	if !config.EnableEagerLoading {
		t.Error("Expected eager loading to be enabled by default")
	}

	if !config.EnableStreaming {
		t.Error("Expected streaming to be enabled by default")
	}

	if config.StreamingThreshold != 100 {
		t.Errorf("Expected streaming threshold to be 100, got %d", config.StreamingThreshold)
	}

	if config.MaxConcurrentExports != 5 {
		t.Errorf("Expected max concurrent exports to be 5, got %d", config.MaxConcurrentExports)
	}

	if config.ExportTimeout != 5*time.Minute {
		t.Errorf("Expected export timeout to be 5 minutes, got %v", config.ExportTimeout)
	}

	if !config.EnableProgressTracking {
		t.Error("Expected progress tracking to be enabled by default")
	}
}

func TestPerformanceOptimizer(t *testing.T) {
	config := DefaultPerformanceConfig()
	config.MaxConcurrentExports = 2 // Limit for testing
	optimizer := NewPerformanceOptimizer(config)

	if optimizer.config.MaxConcurrentExports != 2 {
		t.Errorf("Expected max concurrent exports to be 2, got %d", optimizer.config.MaxConcurrentExports)
	}

	// Test semaphore capacity
	if cap(optimizer.semaphore) != 2 {
		t.Errorf("Expected semaphore capacity to be 2, got %d", cap(optimizer.semaphore))
	}
}

func TestProgressTracking(t *testing.T) {
	optimizer := NewPerformanceOptimizer(DefaultPerformanceConfig())

	// Test progress creation
	collectionID := int64(123)
	format := FormatPostman
	progressID := optimizer.createProgressTracker(collectionID, format)

	if progressID == "" {
		t.Error("Expected non-empty progress ID")
	}

	// Test progress retrieval
	progress, exists := optimizer.GetProgress(progressID)
	if !exists {
		t.Error("Expected progress to exist")
	}

	if progress.CollectionID != collectionID {
		t.Errorf("Expected collection ID %d, got %d", collectionID, progress.CollectionID)
	}

	if progress.Format != string(format) {
		t.Errorf("Expected format %s, got %s", format, progress.Format)
	}

	if progress.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", progress.Status)
	}

	// Test progress update
	optimizer.updateProgress(progressID, "processing", 0.5, "Processing endpoints")

	progress, exists = optimizer.GetProgress(progressID)
	if !exists {
		t.Error("Expected progress to exist after update")
	}

	if progress.Status != "processing" {
		t.Errorf("Expected status 'processing', got %s", progress.Status)
	}

	if progress.Progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", progress.Progress)
	}

	if progress.CurrentStep != "Processing endpoints" {
		t.Errorf("Expected current step 'Processing endpoints', got %s", progress.CurrentStep)
	}

	// Test completion
	optimizer.updateProgress(progressID, "completed", 1.0, "Export completed")

	progress, exists = optimizer.GetProgress(progressID)
	if !exists {
		t.Error("Expected progress to exist after completion")
	}

	if progress.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", progress.Status)
	}

	if progress.CompletedAt == nil {
		t.Error("Expected completed_at to be set")
	}
}

func TestOptimizedExport(t *testing.T) {
	config := DefaultPerformanceConfig()
	config.ExportTimeout = 1 * time.Second // Short timeout for testing
	optimizer := NewPerformanceOptimizer(config)

	// Create test collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: []EndpointWithDetails{}, // Small collection, won't trigger streaming
	}

	// Create mock exporter
	exporter := NewPostmanExporter()
	options := DefaultExportOptions()
	options.Format = FormatPostman

	ctx := context.Background()

	// Test regular export (small collection)
	result, err := optimizer.OptimizedExport(ctx, collection, FormatPostman, options, exporter)
	if err != nil {
		t.Errorf("Optimized export failed: %v", err)
	}

	if result == nil {
		t.Error("Expected export result, got nil")
	}
}

func TestStreamingExport(t *testing.T) {
	config := DefaultPerformanceConfig()
	config.StreamingThreshold = 5 // Low threshold to trigger streaming
	optimizer := NewPerformanceOptimizer(config)

	// Create large test collection
	endpoints := make([]EndpointWithDetails, 10) // Above streaming threshold
	for i := 0; i < 10; i++ {
		endpoints[i] = EndpointWithDetails{
			Endpoint: Endpoint{
				ID:        int64(i + 1),
				Name:      fmt.Sprintf("Endpoint %d", i+1),
				Method:    "GET",
				URL:       fmt.Sprintf("/api/endpoint%d", i+1),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
	}

	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Large Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: endpoints,
	}

	// Create mock exporter
	exporter := NewPostmanExporter()
	options := DefaultExportOptions()
	options.Format = FormatPostman

	ctx := context.Background()

	// Test streaming export (large collection)
	result, err := optimizer.OptimizedExport(ctx, collection, FormatPostman, options, exporter)
	if err != nil {
		t.Errorf("Streaming export failed: %v", err)
	}

	if result == nil {
		t.Error("Expected export result, got nil")
	}
}

func TestConcurrencyControl(t *testing.T) {
	config := DefaultPerformanceConfig()
	config.MaxConcurrentExports = 2 // Limit concurrent exports
	optimizer := NewPerformanceOptimizer(config)

	// Create test collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: []EndpointWithDetails{},
	}

	exporter := NewPostmanExporter()
	options := DefaultExportOptions()
	options.Format = FormatPostman

	// Start multiple exports concurrently
	results := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			ctx := context.Background()
			_, err := optimizer.OptimizedExport(ctx, collection, FormatPostman, options, exporter)
			results <- err
		}()
	}

	// Wait for all exports to complete
	for i := 0; i < 5; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Concurrent export %d failed: %v", i, err)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	optimizer := NewPerformanceOptimizer(DefaultPerformanceConfig())

	// Create test collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: []EndpointWithDetails{},
	}

	exporter := NewPostmanExporter()
	options := DefaultExportOptions()
	options.Format = FormatPostman

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test export with cancelled context
	_, err := optimizer.OptimizedExport(ctx, collection, FormatPostman, options, exporter)
	if err == nil {
		t.Error("Expected error from cancelled context, got nil")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestPerformanceStats(t *testing.T) {
	optimizer := NewPerformanceOptimizer(DefaultPerformanceConfig())

	// Get initial stats
	stats := optimizer.GetStats()
	if stats.ActiveExports != 0 {
		t.Errorf("Expected 0 active exports, got %d", stats.ActiveExports)
	}

	if stats.CompletedExports != 0 {
		t.Errorf("Expected 0 completed exports, got %d", stats.CompletedExports)
	}

	// Create some progress entries manually for testing
	progressID1 := optimizer.createProgressTracker(1, FormatPostman)
	progressID2 := optimizer.createProgressTracker(2, FormatOpenAPI)

	optimizer.updateProgress(progressID1, "processing", 0.5, "Processing")
	optimizer.updateProgress(progressID2, "completed", 1.0, "Completed")

	// Get updated stats
	stats = optimizer.GetStats()
	if stats.ActiveExports != 1 {
		t.Errorf("Expected 1 active export, got %d", stats.ActiveExports)
	}

	if stats.CompletedExports != 1 {
		t.Errorf("Expected 1 completed export, got %d", stats.CompletedExports)
	}
}

func TestGetAllProgress(t *testing.T) {
	optimizer := NewPerformanceOptimizer(DefaultPerformanceConfig())

	// Create multiple progress trackers
	progressID1 := optimizer.createProgressTracker(1, FormatPostman)
	progressID2 := optimizer.createProgressTracker(2, FormatOpenAPI)
	progressID3 := optimizer.createProgressTracker(3, FormatInsomnia)

	optimizer.updateProgress(progressID1, "processing", 0.3, "Processing 1")
	optimizer.updateProgress(progressID2, "processing", 0.7, "Processing 2")
	optimizer.updateProgress(progressID3, "completed", 1.0, "Completed")

	// Get all progress
	allProgress := optimizer.GetAllProgress()

	if len(allProgress) != 3 {
		t.Errorf("Expected 3 progress entries, got %d", len(allProgress))
	}

	// Verify each progress entry
	for progressID, progress := range allProgress {
		if progress.ID != progressID {
			t.Errorf("Progress ID mismatch: expected %s, got %s", progressID, progress.ID)
		}

		if progress.CollectionID == 0 {
			t.Error("Expected non-zero collection ID")
		}

		if progress.Format == "" {
			t.Error("Expected non-empty format")
		}
	}
}
