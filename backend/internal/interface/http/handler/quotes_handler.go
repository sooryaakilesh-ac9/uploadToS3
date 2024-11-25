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
	"os"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Quote represents a single quote with metadata
type Quote struct {
	ID   int      `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Lang string   `json:"lang"`
}

// Schema defines the format of the quotes data
type Schema struct {
	Format   string `json:"format"`
	Encoding string `json:"encoding"`
	FileType string `json:"fileType"`
}

// QuotesMetadata contains metadata about the quotes collection
type QuotesMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalQuotes int    `json:"totalQuotes"`
	Url         string `json:"url"`
	Schema      Schema `json:"schema"`
}

// Quotes is the top-level structure containing both quotes and metadata
type Quotes struct {
	Quotes   []Quote        `json:"quotes"`
	Metadata QuotesMetadata `json:"metadata"`
}

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

// ReadData reads quote data from a Google Sheet and processes it
func ReadData(service *sheets.Service, spreadsheetID string) error {
	const readRange = "English"

	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to read data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return fmt.Errorf("no data found in sheet")
	}

	quotes := processRows(resp.Values)
	metadata := createMetadata(len(quotes))

	outputData := Quotes{
		Quotes:   quotes,
		Metadata: metadata,
	}

	if err := WriteJSONToFile("quotes_output.json", outputData); err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}

	fmt.Println("JSON data successfully written to quotes_output.json")
	return nil
}

// processRows converts sheet rows into Quote structs
func processRows(rows [][]interface{}) []Quote {
	var quotes []Quote

	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue // Skip header row and invalid rows
		}

		quote := Quote{
			ID:   i,
			Text: fmt.Sprintf("%v", row[1]), // Safely convert interface{} to string
			Tags: processTags(fmt.Sprintf("%v", row[0])),
			Lang: "en-US",
		}

		quotes = append(quotes, quote)
	}

	return quotes
}

// processTags cleans and splits tag string into slice
func processTags(rawTags string) []string {
	cleaned := strings.ReplaceAll(rawTags, " ", "")
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, ",")
}

// createMetadata generates metadata for the quotes collection
func createMetadata(totalQuotes int) QuotesMetadata {
	return QuotesMetadata{
		Version:     "1.0",
		LastUpdated: time.Now().Format(time.RFC3339),
		TotalQuotes: totalQuotes,
		Url:         "https://example.com/quotes", // Update with actual URL
		Schema: Schema{
			Format:   "JSON",
			Encoding: "UTF-8",
			FileType: "text",
		},
	}
}

// WriteJSONToFile saves the quotes data to a JSON file
func WriteJSONToFile(filename string, data Quotes) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
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
	if err := ReadData(service, spreadsheetID); err != nil {
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
	// todo write in s3 bucket(local stack)
}
