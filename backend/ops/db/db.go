package db

import (
	"backend/pkg/quotes"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	// Database connection string
	dsn := "host=localhost user=postgres password=toor dbname=postgres port=5432 sslmode=disable"

	// Open a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Check if the 'quotes' table exists
	var count int64
	if err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_name = 'quotes'").Scan(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}

	// If the 'quotes' table doesn't exist, create it
	if count == 0 {
		log.Println("Table 'quotes' does not exist, creating it...")
		if err := db.AutoMigrate(&quotes.Quote{}); err != nil {
			log.Printf("Failed to migrate: %v", err)
			return nil, err
		}
	} else {
		log.Println("Table 'quotes' already exists, skipping creation.")
	}

	// Test the connection by executing a simple query
	err = db.Exec("SELECT 1").Error
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// Successfully connected
	fmt.Println("Successfully connected to the PostgreSQL database!")

	// Return the DB instance for further usage
	return db, nil
}
