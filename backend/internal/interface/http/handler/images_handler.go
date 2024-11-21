package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	io.Copy(destFile, file)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image %s uploaded successfully\n", handler.Filename)

	// make json file

	// todo write to DB
	// write in s3 bucket
}
