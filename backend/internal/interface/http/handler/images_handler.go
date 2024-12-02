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

// handles import of data(from folder)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success_count": successCount,
		"failure_count": failureCount,
		"total": successCount + failureCount,
	})
}

func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

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

	id, err := dbInsertImage(dbConn, flyer)
	if err != nil {
		return fmt.Errorf("failed to insert into db: %w", err)
	}

	if err := utils.UploadToS3LSImages(sourcePath, filename+fmt.Sprintf("_%v", id)); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
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

	// Convert to JSON (Flyer object)
	flyer, err := convertToJson(filePath, handler.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process image: %v", err), http.StatusInternalServerError)
		return
	}

	// Get database connection using the handler's db field
	dbConn, err := h.db.GetConnection()
	if err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}

	id, err := dbInsertImage(dbConn, flyer)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert image: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the imagesMetadata.json file
	images, err := db.FetchAllImagesFromDB()
	if err != nil {
		log.Printf("Error fetching images from DB: %v", err)
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
		return
	}

	// send to ImagesToJson
	if err := utils.ImagesToJson(images); err != nil {
		log.Printf("Error creating images metadata JSON: %v", err)
		http.Error(w, "Failed to update images metadata", http.StatusInternalServerError)
		return
	}

	// Upload image to LocalStack S3
	if err := utils.UploadToS3LSImages(filePath, handler.Filename+fmt.Sprintf("_%v", id)); err != nil {
		log.Printf("Error uploading to S3: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload image to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image %s uploaded successfully\n", handler.Filename)
}

func dbInsertImage(dbConn *gorm.DB, flyer images.Flyer) (uint, error) {
	result := dbConn.Create(&flyer)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to insert image: %w", result.Error)
	}
	log.Printf("Inserted image: %+v", flyer)
	return flyer.Id, nil
}

func convertToJson(filePath, filename string) (images.Flyer, error) {
	// Reopen the file for image analysis
	imgFile, err := os.Open(filePath)
	if err != nil {
		return images.Flyer{}, fmt.Errorf("failed to read uploaded image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	img, format, err := image.Decode(imgFile)
	if err != nil {
		return images.Flyer{}, fmt.Errorf("failed to decode image: %v", err)
	}

	// Validate image format
	validFormats := map[string]bool{"jpeg": true, "png": true, "gif": true}
	if !validFormats[format] {
		return images.Flyer{}, fmt.Errorf("invalid image format")
	}

	// Get the width and height of the image
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Determine file format (uppercase)
	fileFormat := strings.ToUpper(format)
	if fileFormat == "JPG" {
		fileFormat = "JPEG"
	}

	// Determine orientation
	orientation := "portrait"
	if width > height {
		orientation = "landscape"
	}

	// Create Flyer object
	flyer := images.Flyer{
		Design: images.Design{
			TemplateId: "", // Add template ID logic if applicable
			Resolution: images.Resolution{
				Width:  width,
				Height: height,
				Unit:   1, // Assuming pixels
			},
			Type:        "image",
			Tags:        extractImageTags(filename),
			FileFormat:  fileFormat,
			Orientation: orientation,
		},
		Lang: "en-US", // Implement language detection
		Url:  filePath,
	}

	return flyer, nil
}

// Helper function to extract tags from filename
func extractImageTags(filename string) []string {
	// Basic tag extraction - implement more sophisticated logic as needed
	tags := []string{
		strings.TrimSuffix(filename, filepath.Ext(filename)),
		filepath.Ext(filename)[1:],
	}
	return tags
}
