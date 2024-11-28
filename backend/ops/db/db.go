package db

import (
	"backend/pkg/images"
	"backend/pkg/quotes"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// todo implement unit test cases
// todo add the hardcoded values to the ENV file

// connect to DB and returns an instance of the DB
// ConnectToDB connects to the PostgreSQL database and ensures the necessary tables exist
func ConnectToDB() (*gorm.DB, error) {
	// Database connection string (todo: move to environment variables)
	dsn := "host=localhost user=postgres password=toor dbname=postgres port=5432 sslmode=disable"

	// Open a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Migrate the 'quotes' table (if it doesn't exist, AutoMigrate will create it)
	if err := db.AutoMigrate(&quotes.Quote{}); err != nil {
		log.Printf("Failed to migrate 'quotes' table: %v", err)
		return nil, err
	} else {
		log.Println("Table 'quotes' migration successful or already exists.")
	}

	// Migrate the 'flyers' table (if it doesn't exist, AutoMigrate will create it)
	if err := db.AutoMigrate(&images.Flyer{}); err != nil {
		log.Printf("Failed to migrate 'flyers' table: %v", err)
		return nil, err
	} else {
		log.Println("Table 'flyers' migration successful or already exists.")
	}

	// Test the connection by executing a simple query
	if err := db.Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// Successfully connected
	fmt.Println("Successfully connected to the PostgreSQL database!")

	// Return the DB instance for further usage
	return db, nil
}
