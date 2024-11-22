package handler

import (
	"backend/pkg/quotes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func extractSpreadsheetID(sheetLink string) (string, error) {
	// Parse the URL to validate it and extract components
	parsedURL, err := url.Parse(sheetLink)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	// Check if the URL contains `/d/` which is used for spreadsheet IDs
	parts := strings.Split(parsedURL.Path, "/")
	for i, part := range parts {
		if part == "d" && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	// Return an error if the ID is not found
	return "", fmt.Errorf("spreadsheet ID not found in URL")
}

func getService() (*sheets.Service, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("/Users/sooryaakilesh/Downloads/contentservice-442500-a653dca5bcda.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to create Sheets service: %v", err)
	}
	return srv, nil
}

func readData(service *sheets.Service, spreadsheetID string) error {
	// controls what should be read from which sheet
	readRange := "English"

	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to read data: %v", err)
	}

	// Check if there are values in the sheet
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return nil
	}

	// Print each cell with row and column indices
	for rowIndex, row := range resp.Values {
		for colIndex, cell := range row {
			fmt.Printf("Row %d, Column %d: %v\n", rowIndex+1, colIndex+1, cell)
		}
	}
	return nil
}

// handles import of data in the form of excel of a google sheet
func HandleQuotesImport(w http.ResponseWriter, r *http.Request) {
	var v quotes.GoogleSheetsLink

	// Read the body of the request
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON data into the `url` struct
	if err := json.Unmarshal(data, &v); err != nil {
		http.Error(w, "Unable to unmarshal JSON data", http.StatusBadRequest)
		return
	}

	// Extract the spreadsheet ID from the Google Sheets link
	spreadsheetID, err := extractSpreadsheetID(v.GoogleSheetsLink)
	if err != nil {
		http.Error(w, "Invalid Google Sheets link", http.StatusBadRequest)
		return
	}

	fmt.Printf("Extracted Spreadsheet ID: %v\n", spreadsheetID)
	// Create the Sheets service
	service, err := getService()
	if err != nil {
		log.Fatalf("Failed to create Sheets service: %v", err)
	}

	// Read data from the sheet
	if err := readData(service, spreadsheetID); err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}

	// TODO: Implement logic to write to the database

	// TODO: Implement logic to upload to the S3 bucket

	// Respond with a success message (you may want to adjust this based on actual implementation)
	w.WriteHeader(http.StatusOK)
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
