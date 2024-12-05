package handler

import (
	"backend/ops/db"
	"backend/pkg/images"
	"backend/utils"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type ImageHandler struct {
	db db.Database
}

func NewImageHandler(database db.Database) *ImageHandler {
	return &ImageHandler{db: database}
}

// Handles import of data (from folder)
func (h *ImageHandler) HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	importDir := os.Getenv("IMPORT_DIR_IMAGES")
	if importDir == "" {
		http.Error(w, "IMPORT_DIR_IMAGES environment variable not set", http.StatusBadRequest)
		return
	}

	// Validate import directory exists
	if _, err := os.Stat(importDir); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("Import directory %s does not exist", importDir), http.StatusBadRequest)
		return
	}

	// Process each file in the directory
	entries, err := os.ReadDir(importDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read directory: %v", err), http.StatusInternalServerError)
		return
	}

	var successCount, failureCount int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !isValidImageFile(filename) {
			failureCount++
			continue
		}

		if err := h.processImageFile(importDir, filename); err != nil {
			log.Printf("Failed to process %s: %v", filename, err)
			failureCount++
			continue
		}
		successCount++
	}

	// Update the imagesMetadata.json file
	dbConn, err := h.db.GetConnection()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		http.Error(w, "unable to connect to DB", http.StatusInternalServerError)
		return
	}

	var imagesList []images.Flyer
	result := dbConn.Find(&imagesList)
	if result.Error != nil {
		log.Printf("Error fetching images from DB: %v", result.Error)
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
		return
	}

	// Generate images metadata JSON
	flyers := images.Flyers{Flyers: imagesList} // Wrap the images in Flyers
	if err := utils.ImagesToJson(flyers); err != nil {
		log.Printf("Error creating images metadata JSON: %v", err)
		http.Error(w, "Failed to update images metadata", http.StatusInternalServerError)
		return
	}

	// Upload imagesMetadata.json to S3
	metadataFilePath := "./cmd/server/imagesMetadata.json"
	if err := utils.UploadImagesMetadataToS3(metadataFilePath, "imagesMetadata.json"); err != nil {
		log.Printf("Error uploading imagesMetadata.json to S3: %v", err)
		http.Error(w, "unable to update imagesMetadata to S3 bucket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success_count": successCount,
		"failure_count": failureCount,
		"total":         successCount + failureCount,
	})
}

// Checks if a file is a valid image
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

// Processes a single image file
func (h *ImageHandler) processImageFile(importDir, filename string) error {
	sourcePath := filepath.Join(importDir, filename)
	flyer, err := convertToJson(sourcePath, filename)
	if err != nil {
		return fmt.Errorf("failed to convert to flyer: %w", err)
	}

	dbConn, err := h.db.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}

	id, err := dbInsertImage(dbConn, &flyer)
	if err != nil {
		return fmt.Errorf("failed to insert into db: %w", err)
	}

	// Update the URL to include the S3 path and ID
	flyer.Url = fmt.Sprintf("%s/%d_%s", os.Getenv("S3_BUCKET_NAME"), id, filename)
	if err := dbConn.Save(&flyer).Error; err != nil {
		return fmt.Errorf("failed to update flyer URL in db: %w", err)
	}

	if err := utils.UploadToS3Images(sourcePath, fmt.Sprintf("%d_%s", id, filename)); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

// Inserts an image record into the database
func dbInsertImage(dbConn *gorm.DB, flyer *images.Flyer) (uint, error) {
	result := dbConn.Create(flyer)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to insert image: %w", result.Error)
	}
	log.Printf("Inserted image: %+v", flyer)
	return flyer.Id, nil
}

// Converts an image file into a Flyer object
func convertToJson(filePath, filename string) (images.Flyer, error) {
	S3BUCKET_NAME := os.Getenv("S3_BUCKET_NAME")

	imgFile, err := os.Open(filePath)
	if err != nil {
		return images.Flyer{}, fmt.Errorf("failed to read uploaded image: %v", err)
	}
	defer imgFile.Close()

	img, format, err := image.Decode(imgFile)
	if err != nil {
		return images.Flyer{}, fmt.Errorf("failed to decode image: %v", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fileFormat := strings.ToUpper(format)
	if fileFormat == "JPG" {
		fileFormat = "JPEG"
	}

	orientation := "portrait"
	if width > height {
		orientation = "landscape"
	}

	flyer := images.Flyer{
		Design: images.Design{
			TemplateId: "", // Add template logic if needed
			Resolution: images.Resolution{
				Width:  width,
				Height: height,
				Unit:   1, // Assuming pixels
			},
			Type:        "image",
			Tags:        extractImageTags(filename),
			FileFormat:  fileFormat,
			Orientation: orientation,
			FileName: filename,
		},
		Lang: "en-US",
		Url:  fmt.Sprintf("%s/%s", S3BUCKET_NAME, filename), // Temporary URL before DB ID is added
	}

	return flyer, nil
}

// Extracts tags from an image filename
func extractImageTags(filename string) []string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return strings.Split(name, "_")
}

// handles the upload of a single image
func (h *ImageHandler) HandleImagesUpload(w http.ResponseWriter, r *http.Request) {
	// Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the images directory if it doesn't exist
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		os.Mkdir("images", os.ModePerm)
	}

	// Save the file to the images directory
	filePath := filepath.Join("images", handler.Filename)
	destFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy file content
	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Create a Flyer object
	flyer := images.Flyer{
		Url: fmt.Sprintf("%s/%s", os.Getenv("S3_BUCKET_NAME"), handler.Filename), // Temporary URL
		// Add other fields as necessary
	}

	// Save flyer to the database
	dbConn, err := h.db.GetConnection()
	if err != nil {
		http.Error(w, "Failed to get DB connection", http.StatusInternalServerError)
		return
	}

	id, err := dbInsertImage(dbConn, &flyer) // Insert and get ID
	if err != nil {
		http.Error(w, "Failed to insert image into DB", http.StatusInternalServerError)
		return
	}

	// Update the URL to include the ID
	flyer.Url = fmt.Sprintf("%s/%d_%s", os.Getenv("S3_BUCKET_NAME"), id, handler.Filename)

	// Update the flyer in the database with the new URL
	if err := dbConn.Save(&flyer).Error; err != nil {
		http.Error(w, "Failed to update flyer URL in DB", http.StatusInternalServerError)
		return
	}

	var imagesList []images.Flyer
	result := dbConn.Find(&imagesList)
	if result.Error != nil {
		log.Printf("Error fetching images from DB: %v", result.Error)
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
		return
	}

	// Generate images metadata JSON
	flyers := images.Flyers{Flyers: imagesList} // Wrap the images in Flyers
	if err := utils.ImagesToJson(flyers); err != nil {
		log.Printf("Error creating images metadata JSON: %v", err)
		http.Error(w, "Failed to update images metadata", http.StatusInternalServerError)
		return
	}

	// Upload imagesMetadata.json to S3
	metadataFilePath := "./cmd/server/imagesMetadata.json"
	if err := utils.UploadImagesMetadataToS3(metadataFilePath, "imagesMetadata.json"); err != nil {
		log.Printf("Error uploading imagesMetadata.json to S3: %v", err)
		http.Error(w, "unable to update imagesMetadata to S3 bucket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Image uploaded successfully",
		"filename": handler.Filename,
	})
}
