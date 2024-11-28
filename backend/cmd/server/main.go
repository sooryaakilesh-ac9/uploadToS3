package main

import (
	"backend/internal/interface/http/router"
	"backend/ops/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	// Get the directory of the current file
	_, currentFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(currentFile), "..", "..")
	envPath := filepath.Join(rootDir, ".env")
	
	// Load environment variables from .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file at %s: %v", envPath, err)
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
		log.Printf("No PORT specified in .env, using default: %s", PORT)
	}
	fmt.Printf("%v", PORT)

	mux := http.NewServeMux()
	router.RegisterHandlers(mux)

	// test connection ping
	_, err := db.ConnectToDB()
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), mux); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
