package db

import (
	"backend/pkg/quotes"
	"fmt"

	"gorm.io/gorm"
)

// FetchQuoteFromDB fetches a quote by its ID from the database
func FetchQuoteFromDB(quoteId int) (*quotes.Quote, error) {
	// Connect to the database
	db, err := ConnectToDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Declare a variable to hold the fetched quote
	var quote quotes.Quote

	// Fetch the quote from the database by ID
	if err := db.First(&quote, quoteId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If no record is found, return a custom error message
			return nil, fmt.Errorf("quote with ID %d not found", quoteId)
		}
		// If there's a different error, return it
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	// Return the fetched quote
	return &quote, nil
}
