package handler

import (
	"backend/pkg/images"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// handles import of data
func HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	// todo verify the json before handling (use middleware)

	// todo write to DB
	// write in s3 bucket
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

	// Print flyer metadata for debugging purposes
	fmt.Printf("%+v", flyer)

	// TODO: Persist flyer metadata to database
	// saveFlyerToDatabase(flyer)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image %s uploaded successfully\n", handler.Filename)
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
		Id: generateUniqueID(), // You'll need to implement this function
		Design: images.Design{
			TemplateId: "", // Add template ID logic if applicable
			Resolution: images.Resolution{
				Width:  width,
				Height: height,
				Unit:   1, // Assuming pixels
			},
			Type:        "image",
			Tags:        extractImageTags(filename), // Implement tag extraction
			FileFormat:  fileFormat,
			Orientation: orientation,
		},
		Lang: "en-US", // Implement language detection
		Url:  filePath,
	}

	return flyer, nil
}

// Helper function to generate unique ID
func generateUniqueID() string {
	return fmt.Sprintf("flyer_%d", time.Now().UnixNano())
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
