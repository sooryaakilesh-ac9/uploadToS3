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

// handles import of data(from folder)
func HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	// Specify the directory to import images from
	importDir := "/Users/sooryaakilesh/Documents/contentService/designs"

	// Validate import directory exists
	if _, err := os.Stat(importDir); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("Import directory %s does not exist", importDir), http.StatusBadRequest)
		return
	}

	// Create images directory if it doesn't exist
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		if err := os.Mkdir("images", os.ModePerm); err != nil {
			http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
			return
		}
	}

	// Connect to database
	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	// Track import statistics
	var successCount, failureCount int
	var importedImages []images.Flyer

	// Read directory contents
	entries, err := os.ReadDir(importDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read import directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Process each file in the import directory
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Check file extension
		filename := entry.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			log.Printf("Skipping non-image file: %s", filename)
			failureCount++
			continue
		}

		// Construct full source and destination paths
		sourcePath := filepath.Join(importDir, filename)
		destPath := filepath.Join("images", filename)

		// Copy file to images directory
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			log.Printf("Failed to open source file %s: %v", sourcePath, err)
			failureCount++
			continue
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			log.Printf("Failed to create destination file %s: %v", destPath, err)
			failureCount++
			continue
		}
		defer destFile.Close()

		// Copy file contents
		if _, err := io.Copy(destFile, sourceFile); err != nil {
			log.Printf("Failed to copy file %s: %v", filename, err)
			failureCount++
			continue
		}

		// Convert to Flyer object
		flyer, err := convertToJson(destPath, filename)
		if err != nil {
			log.Printf("Failed to convert %s to Flyer: %v", filename, err)
			failureCount++
			continue
		}

		// Insert into database
		if err := dbInsertImage(dbConn, flyer); err != nil {
			log.Printf("Failed to insert image %s into database: %v", filename, err)
			failureCount++
			continue
		}

		// Upload to S3
		if err := utils.UploadToS3LSImages(destPath, filename); err != nil {
			log.Printf("Failed to upload %s to S3: %v", filename, err)
			failureCount++
			continue
		}

		// Track successful import
		successCount++
		importedImages = append(importedImages, flyer)
	}

	// Update images metadata JSON
	if err := utils.ImagesToJson(importedImages); err != nil {
		log.Printf("Failed to update images metadata: %v", err)
	}

	// Prepare and send response
	response := struct {
		Message       string `json:"message"`
		SuccessCount  int    `json:"success_count"`
		FailureCount  int    `json:"failure_count"`
		TotalAttempts int    `json:"total_attempts"`
	}{
		Message:       "Image import completed",
		SuccessCount:  successCount,
		FailureCount:  failureCount,
		TotalAttempts: successCount + failureCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handles the upload of a single image
func HandleImagesUpload(w http.ResponseWriter, r *http.Request) {
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

	// Write to database
	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}

	if err := dbInsertImage(dbConn, flyer); err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert image: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch image from the database
	image, err := db.FetchImageFromDB(1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch image from DB: %v", err), http.StatusInternalServerError)
		return
	}

	// Print image JSON (for debugging purposes)
	fmt.Printf("image JSON => %v\n", image)

	// Update the imagesMetadata.json file
	images, err := db.FetchAllImagesFromDB()
	if err != nil {
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
		return
	}

	// send to quotesToJson
	if err := utils.ImagesToJson(images); err != nil {
		http.Error(w, "Failed to update images metadata", http.StatusInternalServerError)
		return
	}

	// Upload image to LocalStack S3
	if err := utils.UploadToS3LSImages(filePath, handler.Filename); err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload image to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image %s uploaded successfully\n", handler.Filename)
}

func dbInsertImage(dbConn *gorm.DB, flyer images.Flyer) error {
	result := dbConn.Create(&flyer)
	if result.Error != nil {
		return fmt.Errorf("failed to insert image: %w", result.Error)
	}
	log.Printf("Inserted image: %+v", flyer)
	return nil
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