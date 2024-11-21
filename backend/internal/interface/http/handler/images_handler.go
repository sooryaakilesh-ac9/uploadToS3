package handler

import "net/http"

// handles import of data
func HandleImagesImport(w http.ResponseWriter, r *http.Request) {
	// todo verify the json before handling (use middleware)

	// todo write to DB
	// write in s3 bucket
}

// handles the upload of a single image
func HandleImagesUpload(w http.ResponseWriter, r *http.Request) {
	// todo verify the json before handling (use middleware)

	// todo write to DB
	// write in s3 bucket
}
