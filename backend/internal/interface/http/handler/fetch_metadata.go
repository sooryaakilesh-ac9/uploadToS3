package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ImageMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	Total       int    `json:"total"`
	Url         string `json:"url"`
}

type QuoteMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	Total       int    `json:"total"`
	Url         string `json:"url"`
}

type CombinedMetadata struct {
	Quote *QuoteMetadata `json:"quote,omitempty"`
	Image *ImageMetadata `json:"image,omitempty"`
}

// Unified MetadataLoader to reduce redundancy
func loadMetadata(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return nil
}

// HandleFetchMetadata refactored
func HandleFetchMetadata(w http.ResponseWriter, r *http.Request) {
	// Validate path
	if !strings.HasPrefix(r.URL.Path, "/metadata/media") {
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	// Extract media types
	queryParams := r.URL.Query()
	mediaTypes := queryParams["type"]
	if len(mediaTypes) == 0 {
		mediaTypes = []string{"quote", "image"}
	}

	// Validate media types
	for _, mediaType := range mediaTypes {
		if mediaType != "quote" && mediaType != "image" {
			http.Error(w, "Invalid 'type' value. Allowed values are 'quote' or 'image'", http.StatusBadRequest)
			return
		}
	}

	// Prepare combined metadata
	var combinedMetadata CombinedMetadata

	for _, mediaType := range mediaTypes {
		switch mediaType {
		case "quote":
			filePath := "/Users/sooryaakilesh/Documents/contentService/backend/cmd/server/quotesMetadata.json"
			var jsonData MetaDataQuotes
			if err := loadMetadata(filePath, &jsonData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			combinedMetadata.Quote = &QuoteMetadata{
				Version:     jsonData.QuotesMetadata.Version,
				LastUpdated: jsonData.QuotesMetadata.LastUpdated,
				Total:       jsonData.QuotesMetadata.TotalQuotes,
				Url:         jsonData.QuotesMetadata.URL,
			}
		case "image":
			filePath := "/Users/sooryaakilesh/Documents/contentService/backend/cmd/server/imagesMetadata.json"
			var jsonData MetadataFlyers
			if err := loadMetadata(filePath, &jsonData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			combinedMetadata.Image = &ImageMetadata{
				Version:     jsonData.FlyersMetadata.Version,
				LastUpdated: jsonData.FlyersMetadata.LastUpdated,
				Total:       jsonData.FlyersMetadata.TotalFlyers,
				Url:         jsonData.FlyersMetadata.Url,
			}
		}
	}

	// Construct and send response
	response, err := json.Marshal(combinedMetadata)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
