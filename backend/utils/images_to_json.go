package utils

import (
	"backend/pkg/images"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MetadataFlyers struct {
	Flyers         []Flyer        `json:"flyers"`
	FlyersMetadata FlyersMetadata `json:"metadata"`
}

type Flyer struct {
	Id     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Design Design `json:"design" gorm:"embedded"`
	Lang   string `json:"lang"`
	Url    string `json:"url"`
}

type Design struct {
	TemplateId  string     `json:"templateId" gorm:"column:template_id"`
	Resolution  Resolution `json:"resolution" gorm:"embedded"`
	Type        string     `json:"type"`
	Tags        []string   `json:"tags" gorm:"serializer:json"`
	FileFormat  string     `json:"fileFormat" gorm:"column:file_format"`
	Orientation string     `json:"orientation"`
	FileName    string     `json:"fileName"`
}

type Resolution struct {
	Width  int `json:"width" gorm:"column:width"`
	Height int `json:"height" gorm:"column:height"`
	Unit   int `json:"unit" gorm:"column:unit"`
}

type FlyersMetadata struct {
	LastUpdated string `json:"lastUpdated"`
	TotalFlyers int    `json:"total"`
	Url         string `json:"url"`
	Version     uint   `json:"version"`
}

type MetaDataImages struct {
	Images         []images.Flyer `json:"media"`
	ImagesMetadata FlyersMetadata `json:"metadata"`
}

// ImagesToJson converts the flyer images to a JSON file with metadata
func ImagesToJson(images images.Flyers) error {
	metadataPath := os.Getenv("IMAGE_METADATA_PATH")
	metadataFileName := os.Getenv("IMAGE_METADATA_FILENAME")

	// Ensure the directory exists
	if err := os.MkdirAll(metadataPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	metadata := FlyersMetadata{
		Version:     1,
		LastUpdated: time.Now().Format(time.RFC3339),
		TotalFlyers: len(images.Flyers),
		Url:         os.Getenv("IMAGE_METADATA_URL"), // Replace with actual URL or dynamic generation
	}

	imageData := MetaDataImages{
		Images:         images.Flyers,
		ImagesMetadata: metadata,
	}

	// Marshal the data into JSON with indentation for readability
	imageJson, err := json.MarshalIndent(imageData, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to convert images to JSON: %v", err)
	}

	// Create the JSON file to store the result
	filePath := filepath.Join(metadataPath, metadataFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create JSON file: %v", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(imageJson)
	if err != nil {
		return fmt.Errorf("unable to write JSON data to file: %v", err)
	}

	UploadImagesMetadataToS3(metadataPath, metadataFileName)

	fmt.Printf("JSON file %v has been created successfully.\n", metadataFileName)
	return nil
}
