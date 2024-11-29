package db

import (
	"backend/pkg/quotes"
	"fmt"
	"log"
)

func InsertQuoteToDB(quote quotes.Quote) error {
	// Connect to the database
	db, err := GetDB()
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Ensure the table exists with correct schema
	err = db.AutoMigrate(&quotes.Quote{})
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// Insert the quote into the database using SQL directly to debug
	result := db.Debug().Create(quote) // Added Debug() to see the SQL query
	if result.Error != nil {
		return fmt.Errorf("failed to insert quote into the database: %w", result.Error)
	}

	// Log the inserted quote ID
	log.Printf("Inserted quote with ID: %d", quote.Id)

	return nil
}
