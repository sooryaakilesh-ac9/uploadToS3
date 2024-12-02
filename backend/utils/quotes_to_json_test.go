package utils

import (
	"backend/pkg/quotes"
	"os"
	"path/filepath"
	"testing"
)

func TestQuotesToJson(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	// tmpFile := filepath.Join(tmpDir, "test_quotes_metadata.json")
	
	// Set all required environment variables
	os.Setenv("QUOTE_METADATA_FILENAME", "quotesMetadata.json")
	os.Setenv("QUOTE_METADATA_PATH", tmpDir)
	os.Setenv("S3_REGION", "us-east-1")
	os.Setenv("S3_ENDPOINT", "http://localhost:4566")
	os.Setenv("S3_ID", "test")
	os.Setenv("S3_SECRET", "test")
	os.Setenv("S3_TOKEN", "")
	os.Setenv("S3_BUCKET_NAME", "test-bucket")
	os.Setenv("S3_QUOTES_DIR_PATH", "/quotes/")

	// Store original env values
	originalEnv := map[string]string{
		"QUOTE_METADATA_FILENAME": os.Getenv("QUOTE_METADATA_FILENAME"),
		"QUOTE_METADATA_PATH":     os.Getenv("QUOTE_METADATA_PATH"),
		"S3_REGION":              os.Getenv("S3_REGION"),
		"S3_ENDPOINT":            os.Getenv("S3_ENDPOINT"),
		"S3_ID":                  os.Getenv("S3_ID"),
		"S3_SECRET":              os.Getenv("S3_SECRET"),
		"S3_TOKEN":               os.Getenv("S3_TOKEN"),
		"S3_BUCKET_NAME":         os.Getenv("S3_BUCKET_NAME"),
		"S3_QUOTES_DIR_PATH":     os.Getenv("S3_QUOTES_DIR_PATH"),
	}

	// Restore original environment variables after test
	defer func() {
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	tests := []struct {
		name    string
		quotes  []quotes.Quote
		wantErr bool
	}{
		{
			name:    "empty quotes list",
			quotes:  []quotes.Quote{},
			wantErr: false,
		},
		{
			name: "valid quotes list",
			quotes: []quotes.Quote{
				{
					Id:   1,
					Text: "Test quote",
					Tags: []string{"test", "quote"},
					Lang: "en-US",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple quotes",
			quotes: []quotes.Quote{
				{
					Id:   1,
					Text: "First quote",
					Tags: []string{"test", "first"},
					Lang: "en-US",
				},
				{
					Id:   2,
					Text: "Second quote",
					Tags: []string{"test", "second"},
					Lang: "es-ES",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := QuotesToJson(tt.quotes)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuotesToJson() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file was created
			expectedFile := filepath.Join(tmpDir, "quotesMetadata.json")
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) && !tt.wantErr {
				t.Errorf("QuotesToJson() failed to create file at %s", expectedFile)
			}
		})
	}
} 