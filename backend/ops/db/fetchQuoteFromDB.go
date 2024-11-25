package db

import (
	"backend/pkg/quotes"
	"backend/utils"
	"fmt"

	"gorm.io/gorm"
)

// unit tests
//invalid quoteID
//quoteID not present
//different type

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
			return nil, fmt.Errorf("quote with ID %d not found", quoteId)
		}
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	// a standalone module which conveters the given data into JSON format
	utils.JsonHandler(quote)

	// Return the fetched quote
	return &quote, nil
}
