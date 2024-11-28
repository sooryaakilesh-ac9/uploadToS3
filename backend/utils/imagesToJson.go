package utils

import (
	"backend/pkg/images"
	"encoding/json"
	"fmt"
	"os"
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
}

type Resolution struct {
	Width  int `json:"width" gorm:"column:width"`
	Height int `json:"height" gorm:"column:height"`
	Unit   int `json:"unit" gorm:"column:unit"`
}

type FlyersMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalFlyers int    `json:"totalFlyers"`
	Url         string `json:"url"`
	Schema      Schema `json:"schema"`
}

type MetaDataImages struct {
	Images         []images.Flyer    `json:"images"`
	ImagesMetadata FlyersMetadata    `json:"metadata"`
}

// ImagesToJson converts the flyer images to a JSON file with metadata
func ImagesToJson(images []images.Flyer) error {
	metadata := FlyersMetadata{
		Version:     "1.0",
		LastUpdated: time.Now().Format(time.RFC3339),
		TotalFlyers: len(images),
		Url:         os.Getenv("IMAGE_METADATA_URL"), // Replace with actual URL or dynamic generation
		Schema: Schema{
			Format:   "JSON",
			Encoding: "UTF-8",
			FileType: "application/json",
		},
	}

	imageData := MetaDataImages{
		Images:         images,
		ImagesMetadata: metadata,
	}

	// Marshal the data into JSON with indentation for readability
	imageJson, err := json.MarshalIndent(imageData, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to convert images to JSON: %v", err)
	}

	// Create the JSON file to store the result
	file, err := os.Create(os.Getenv("IMAGE_METADATA_FILENAME"))
	if err != nil {
		return fmt.Errorf("unable to create JSON file: %v", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(imageJson)
	if err != nil {
		return fmt.Errorf("unable to write JSON data to file: %v", err)
	}

	metadataPath := os.Getenv("IMAGE_METADATA_PATH")
	metadataFileName := "/" + os.Getenv("IMAGE_METADATA_FILENAME")
	UploadImagesMetadataToS3LS(metadataPath, metadataFileName)

	fmt.Printf("JSON file %v has been created successfully.\n", metadataFileName)
	return nil
}
