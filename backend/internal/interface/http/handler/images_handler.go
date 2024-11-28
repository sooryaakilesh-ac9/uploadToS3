package handler

import (
	"backend/ops/db"
	"backend/pkg/images"
	"backend/utils"
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

// handles import of data(google drive link or from folder)
func HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	// todo implement batch processing

	// todo write to DB
	// write in s3 bucket
}

// handles the upload of a single image
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
	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "failed to connect to database: %w", http.StatusInternalServerError)
		return
	}

	dbInsertImage(dbConn, flyer)

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
	}

	// send to quotesToJson
	utils.ImagesToJson(images)

	// Upload image to LocalStack S3 (New code)
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
		return fmt.Errorf("failed to insert quote: %w", result.Error)
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
