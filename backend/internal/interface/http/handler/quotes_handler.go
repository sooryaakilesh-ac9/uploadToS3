package handler

import "net/http"

// handles import of data in the form of excel of a google sheet
func HandleQuotesImport(w http.ResponseWriter, r *http.Request) {
	// todo verify the json before handling (use middleware)

	// todo write to DB
	// write in s3 bucket
}

// handles the upload of a single quote
func HandleQuotesUpload(w http.ResponseWriter, r *http.Request) {
	// todo verify the json before handling (use middleware)
	// create a quote object based on the input

	// todo write to DB
	// write in s3 bucket
}
