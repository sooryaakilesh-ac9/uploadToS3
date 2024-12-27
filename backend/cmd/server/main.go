package main

import (
	"backend/internal"
	"backend/internal/config"
	"backend/internal/delivery/http/router"
	"backend/internal/infrastructure/persistence/postgres"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatal(err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize handlers
	imageHandler, err := internal.InitializeImageHandler(db)
	if err != nil {
		log.Fatalf("Failed to initialize image handler: %v", err)
	}

	// Setup router
	mux := http.NewServeMux()
	router.RegisterHandlers(mux, imageHandler)

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loadEnv() error {
	_, currentFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(currentFile), "..", "..")
	envPath := filepath.Join(rootDir, ".env")
	
	return godotenv.Load(envPath)
}
