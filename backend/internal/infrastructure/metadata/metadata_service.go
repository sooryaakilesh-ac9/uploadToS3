package metadata

import (
	"backend/internal/domain/entity"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MetadataService struct {
	s3Service S3Service
}

type Metadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	Total       int    `json:"total"`
	Url         string `json:"url"`
}

type ImageMetadata struct {
	Images   []entity.Flyer `json:"media"`
	Metadata Metadata       `json:"metadata"`
}

type QuoteMetadata struct {
	Quotes   []entity.Quote `json:"quotes"`
	Metadata Metadata       `json:"metadata"`
}

func NewMetadataService(s3Service S3Service) *MetadataService {
	return &MetadataService{
		s3Service: s3Service,
	}
}

func (s *MetadataService) UpdateImageMetadata(images []entity.Flyer) error {
	metadata := Metadata{
		Version:     "1",
		LastUpdated: time.Now().Format(time.RFC3339),
		Total:       len(images),
		Url:         os.Getenv("IMAGE_METADATA_URL"),
	}

	imageData := ImageMetadata{
		Images:   images,
		Metadata: metadata,
	}

	return s.saveAndUploadMetadata(imageData, "IMAGE_METADATA_PATH", "IMAGE_METADATA_FILENAME", "imagesMetadata.json")
}

func (s *MetadataService) UpdateQuoteMetadata(quotes []entity.Quote) error {
	metadata := Metadata{
		Version:     "1",
		LastUpdated: time.Now().Format(time.RFC3339),
		Total:       len(quotes),
		Url:         os.Getenv("QUOTE_METADATA_URL"),
	}

	quoteData := QuoteMetadata{
		Quotes:   quotes,
		Metadata: metadata,
	}

	return s.saveAndUploadMetadata(quoteData, "QUOTE_METADATA_PATH", "QUOTE_METADATA_FILENAME", "quotesMetadata.json")
}

func (s *MetadataService) saveAndUploadMetadata(data interface{}, pathEnv, filenameEnv, defaultFilename string) error {
	metadataPath := os.Getenv(pathEnv)
	metadataFileName := os.Getenv(filenameEnv)
	if metadataFileName == "" {
		metadataFileName = defaultFilename
	}

	if err := os.MkdirAll(metadataPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to convert to JSON: %v", err)
	}

	filePath := filepath.Join(metadataPath, metadataFileName)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("unable to write JSON file: %v", err)
	}

	return s.s3Service.UploadMetadata(filePath, metadataFileName)
} 