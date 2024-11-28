package handler

import (
	"backend/ops/db"
	"backend/pkg/quotes"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lib/pq"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

const (
	credentialsFile = "/Users/sooryaakilesh/Downloads/contentservice-442500-a653dca5bcda.json"
	batchSize       = 100
	readRange       = "English"
)

type Quote struct {
	Id   int            `json:"id" gorm:"primaryKey"`
	Text string         `json:"text"`
	Tags pq.StringArray `gorm:"type:text[];column:tags" json:"tags"`
	Lang string         `json:"lang"`
}

type Schema struct {
	Format   string `json:"format"`
	Encoding string `json:"encoding"`
	FileType string `json:"fileType"`
}

type QuotesMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalQuotes int    `json:"totalQuotes"`
	URL         string `json:"url"`
	Schema      Schema `json:"schema"`
}

type Quotes struct {
	Quotes   []Quote        `json:"quotes"`
	Metadata QuotesMetadata `json:"metadata"`
}

func extractSpreadsheetID(sheetLink string) (string, error) {
	parsedURL, err := url.Parse(sheetLink)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}
	parts := strings.Split(parsedURL.Path, "/")
	for i, part := range parts {
		if part == "d" && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("spreadsheet ID not found in URL")
}

func getService() (*sheets.Service, error) {
	ctx := context.Background()
	if credentialsFile == "" {
		return nil, fmt.Errorf("credentials file path is empty")
	}
	return sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
}

func ReadData(service *sheets.Service, spreadsheetID string) ([]Quote, error) {
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to read data from sheet: %w", err)
	}
	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("no data found in sheet")
	}
	return processRows(resp.Values), nil
}

func processRows(rows [][]interface{}) []Quote {
	quotes := make([]Quote, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}
		quotes = append(quotes, Quote{
			Id:   i,
			Text: fmt.Sprintf("%v", row[1]),
			Tags: processTags(fmt.Sprintf("%v", row[0])),
			Lang: "en-US",
		})
	}
	return quotes
}

func processTags(rawTags string) []string {
	cleaned := strings.ReplaceAll(rawTags, " ", "")
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, ",")
}

func HandleQuotesImport(w http.ResponseWriter, r *http.Request) {
	var payload quotes.GoogleSheetsLink
	if err := parseRequestBody(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	spreadsheetID, err := extractSpreadsheetID(payload.GoogleSheetsLink)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid Google Sheets link: %v", err), http.StatusBadRequest)
		return
	}

	service, err := getService()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize Sheets service: %v", err), http.StatusInternalServerError)
		return
	}

	quotes, err := ReadData(service, spreadsheetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := processQuotesInBatches(quotes); err != nil {
		http.Error(w, fmt.Sprintf("Error processing quotes: %v", err), http.StatusInternalServerError)
		return
	}

	// mechanism to update quotesMetadata.json
	quotesJson, err := db.FetchAllQuotesFromDB()
	if err != nil {
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
	}

	// send to quotesToJson
	utils.QuotesToJson(quotesJson)

	writeResponse(w, http.StatusOK, "Data imported successfully")
}

func processQuotesInBatches(quotes []Quote) error {
	dbConn, err := db.ConnectToDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	totalQuotes := len(quotes)
	for i := 0; i < totalQuotes; i += batchSize {
		end := i + batchSize
		if end > totalQuotes {
			end = totalQuotes
		}
		if err := processBatch(dbConn, quotes[i:end]); err != nil {
			return fmt.Errorf("error processing batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func processBatch(dbConn *gorm.DB, batch []Quote) error {
	for _, quote := range batch {
		if err := dbInsertQuote(dbConn, quote); err != nil {
			return err
		}
	}
	return nil
}

func dbInsertQuote(dbConn *gorm.DB, quote Quote) error {
	result := dbConn.Create(&quote)
	if result.Error != nil {
		return fmt.Errorf("failed to insert quote: %w", result.Error)
	}
	log.Printf("Inserted quote: %+v", quote)
	return nil
}

func parseRequestBody(r *http.Request, v interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func HandleQuotesUpload(w http.ResponseWriter, r *http.Request) {
	var quote Quote

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read JSON data", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &quote); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received quote: %+v", quote)

	// Write to database
	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "failed to connect to database: %w", 500)
		return
	}

	dbInsertQuote(dbConn, quote)

	// update the quotesMetadata.json file
	quotes, err := db.FetchAllQuotesFromDB()
	if err != nil {
		http.Error(w, "unable to fetch data from DB", http.StatusInternalServerError)
	}

	// send to quotesToJson
	utils.QuotesToJson(quotes)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Quote uploaded successfully"))
}

func writeResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
