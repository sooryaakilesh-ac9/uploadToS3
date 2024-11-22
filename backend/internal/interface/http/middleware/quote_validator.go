package middleware

import (
	"backend/pkg/quotes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

func isValidGoogleSheetsURL(url string) bool {
	// Regex pattern to match Google Sheets URLs
	const googleSheetsURLPattern = `^https://docs\.google\.com/spreadsheets/d/[a-zA-Z0-9_-]+(/.*)?$`

	re := regexp.MustCompile(googleSheetsURLPattern)

	// Check if the URL matches the pattern
	return re.MatchString(url)
}

// checks if the given link is a valid google spread sheet link
func CheckQuotesLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var url quotes.GoogleSheetsLink

		// checks if the http method is valid or not
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read from JSON", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(data, &url); err != nil {
			http.Error(w, "unable to unmarshall JSON data", http.StatusInternalServerError)
			return
		}

		if isValidGoogleSheetsURL(url.GoogleSheetsLink) {
			fmt.Printf("valid google sheets link\n")
		} else {
			http.Error(w, "invalid google sheets link", http.StatusBadRequest)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// checks if the JSON is valid and contains all fields
// http methods are checked
func QuotesJsonAndMethodValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var quote quotes.Quote

		// checks if the http method is valid or not
		if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Create a map to store the incoming JSON data
		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Use reflection to compare struct fields with JSON keys
		if err := compareStructWithJSON(quote, jsonMap); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Attempt to unmarshal into the quote struct
		err = json.Unmarshal(data, &quote)
		if err != nil {
			http.Error(w, "JSON does not match quote structure", http.StatusBadRequest)
			return
		}
		// Restore the body for subsequent handlers
		r.Body = io.NopCloser(bytes.NewBuffer(data))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// compareStructWithJSON checks if JSON keys match struct field names
func compareStructWithJSON(obj interface{}, jsonMap map[string]interface{}) error {
	// Get the type of the struct
	t := reflect.TypeOf(obj)

	// If it's a pointer, get the underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Collect struct field names (considering JSON tags)
	structFields := make(map[string]bool)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check JSON tag first
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// Split to handle options like "fieldname,omitempty"
			tagName := strings.Split(jsonTag, ",")[0]
			structFields[tagName] = true
		} else {
			// Use field name if no JSON tag
			structFields[strings.ToLower(field.Name)] = true
		}
	}

	// Check if JSON keys match struct fields
	for key := range jsonMap {
		// Convert to lowercase for case-insensitive comparison
		lowercaseKey := strings.ToLower(key)

		if _, exists := structFields[lowercaseKey]; !exists {
			return fmt.Errorf("unexpected JSON key: %s", key)
		}
	}

	// Check if all required struct fields are present
	for field := range structFields {
		if _, exists := jsonMap[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	return nil
}
