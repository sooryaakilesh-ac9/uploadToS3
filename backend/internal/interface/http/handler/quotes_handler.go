package handler

import (
	"backend/pkg/quotes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// handles import of data in the form of excel of a google sheet
func HandleQuotesImport(w http.ResponseWriter, r *http.Request) {
	
	
	// todo write to DB
	// write in s3 bucket
}

// handles the upload of a single quote
func HandleQuotesUpload(w http.ResponseWriter, r *http.Request) {
	// create a quote object based on the input
	var quote quotes.Quote

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read JSON data", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	if err := json.Unmarshal(data, &quote); err != nil {
		http.Error(w, "unable to read JSON data", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v", string(data))

	// todo write to DB
	// write in s3 bucket
}
