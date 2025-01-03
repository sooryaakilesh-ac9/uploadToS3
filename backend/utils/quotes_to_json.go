package utils

import (
	"backend/pkg/quotes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MetaDataQuotes struct {
	Quotes         []quotes.Quote `json:"quotes"`
	QuotesMetadata QuotesMetadata `json:"metadata"`
}

type QuotesMetadata struct {
	LastUpdated string `json:"lastUpdated"`
	TotalQuotes int    `json:"total"`
	Url         string `json:"url"`
	Version     uint   `json:"version"`
}

func QuotesToJson(quotes []quotes.Quote) error {
	metadataPath := os.Getenv("QUOTE_METADATA_PATH")
	metadataFileName := os.Getenv("QUOTE_METADATA_FILENAME")

	// Ensure the directory exists
	if err := os.MkdirAll(metadataPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write JSON data to a file
	filePath := filepath.Join(metadataPath, metadataFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create JSON file at %s: %v", filePath, err)
	}
	defer file.Close()

	// Prepare metadata
	metadata := QuotesMetadata{
		Version:     1,
		LastUpdated: time.Now().Format(time.RFC3339),
		TotalQuotes: len(quotes),
		Url: "my-s3-bucket/quotes.json",
	}

	quote := MetaDataQuotes{
		Quotes:         quotes,
		QuotesMetadata: metadata,
	}

	quoteJson, err := json.MarshalIndent(quote, "", "  ") // Indented JSON for better readability
	if err != nil {
		return fmt.Errorf("unable to convert quotes to JSON: %v", err)
	}

	_, err = file.Write(quoteJson)
	if err != nil {
		return fmt.Errorf("unable to write JSON data to file: %v", err)
	}

	UploadQuotesMetadataToS3(metadataPath, metadataFileName)

	fmt.Println("JSON file 'quotesMetadata.json' has been created successfully.")
	return nil
}
