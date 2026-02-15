package main

import (
	"context"
	"flag"
	"fmt"
	"gin-scalable-api/pkg/swagger"
	"log"
	"os"
	"time"
)

func main() {
	// Define command-line flags
	generateCmd := flag.Bool("generate", false, "Generate Swagger documentation")
	validateCmd := flag.Bool("validate", false, "Validate Swagger annotations")
	watchCmd := flag.Bool("watch", false, "Watch for changes and regenerate")
	outputDir := flag.String("output", "docs", "Output directory for generated files")
	searchDir := flag.String("dir", "./", "Directory to search for annotations")
	helpCmd := flag.Bool("help", false, "Show help message")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Swagger Documentation CLI\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  # Generate documentation\n")
		fmt.Fprintf(os.Stderr, "  %s -generate\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Validate annotations\n")
		fmt.Fprintf(os.Stderr, "  %s -validate\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Watch for changes\n")
		fmt.Fprintf(os.Stderr, "  %s -watch\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Custom output directory\n")
		fmt.Fprintf(os.Stderr, "  %s -generate -output ./api-docs\n\n", os.Args[0])
	}

	flag.Parse()

	// Show help if requested
	if *helpCmd {
		flag.Usage()
		os.Exit(0)
	}

	// Create generator
	gen := swagger.NewGenerator(&swagger.Config{
		SwagPath: "swag",
	})

	ctx := context.Background()

	// Execute command
	if *validateCmd {
		fmt.Println("Validating Swagger annotations...")
		errors, err := gen.Validate(ctx)
		if err != nil {
			log.Fatalf("Validation failed: %v", err)
		}

		if len(errors) == 0 {
			fmt.Println("✓ All annotations are valid")
			os.Exit(0)
		}

		fmt.Printf("Found %d validation errors:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  %s:%d - %s\n", e.File, e.Line, e.Message)
		}
		os.Exit(1)
	}

	if *generateCmd {
		fmt.Println("Generating Swagger documentation...")
		opts := swagger.GenerateOptions{
			OutputDir: *outputDir,
			SearchDir: *searchDir,
			Exclude:   []string{"vendor", "tmp", "bin"},
		}

		if err := gen.Generate(ctx, opts); err != nil {
			log.Fatalf("Generation failed: %v", err)
		}

		fmt.Println("✓ Documentation generated successfully")
		fmt.Printf("  Output: %s/swagger.json, %s/swagger.yaml\n", *outputDir, *outputDir)
		os.Exit(0)
	}

	if *watchCmd {
		fmt.Println("Starting watch mode...")
		opts := swagger.WatchOptions{
			GenerateOptions: swagger.GenerateOptions{
				OutputDir: *outputDir,
				SearchDir: *searchDir,
				Exclude:   []string{"vendor", "tmp", "bin"},
			},
			Interval:      2 * time.Second,
			DebounceDelay: 500 * time.Millisecond,
		}

		if err := gen.Watch(ctx, opts); err != nil {
			log.Fatalf("Watch failed: %v", err)
		}
		os.Exit(0)
	}

	// No command specified
	flag.Usage()
	os.Exit(1)
}
