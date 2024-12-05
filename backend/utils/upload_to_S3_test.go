package utils

import (
	"os"
	"testing"
)

func TestUploadToS3Images(t *testing.T) {
	// Setup test environment
	originalS3ID := os.Getenv("S3_ID")
	originalS3SECRET := os.Getenv("S3_SECRET")
	originalS3TOKEN := os.Getenv("S3_TOKEN")
	
	defer func() {
		os.Setenv("S3_ID", originalS3ID)
		os.Setenv("S3_SECRET", originalS3SECRET)
		os.Setenv("S3_TOKEN", originalS3TOKEN)
	}()

	tests := []struct {
		name     string
		path     string
		fileName string
		setupEnv func()
		wantErr  bool
	}{
		{
			name:     "missing credentials",
			path:     "testdata/test.jpg",
			fileName: "test.jpg",
			setupEnv: func() {
				os.Setenv("S3_ID", "")
				os.Setenv("S3_SECRET", "")
				os.Setenv("S3_TOKEN", "")
			},
			wantErr: true,
		},
		{
			name:     "invalid file path",
			path:     "nonexistent/path.jpg",
			fileName: "test.jpg",
			setupEnv: func() {
				os.Setenv("S3_ID", "test")
				os.Setenv("S3_SECRET", "test")
				os.Setenv("S3_TOKEN", "test")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			err := UploadToS3Images(tt.path, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadToS3Images() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUploadQuotesMetadataToS3(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		fileName string
		setupEnv func()
		wantErr  bool
	}{
		{
			name:     "missing environment variables",
			path:     "testdata/quotes.json",
			fileName: "quotes.json",
			setupEnv: func() {
				os.Setenv("S3_REGION", "")
				os.Setenv("S3_ENDPOINT", "")
			},
			wantErr: true,
		},
		{
			name:     "invalid file path",
			path:     "nonexistent/quotes.json",
			fileName: "quotes.json",
			setupEnv: func() {
				os.Setenv("S3_REGION", "us-east-1")
				os.Setenv("S3_ENDPOINT", "http://localhost:4566")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			err := UploadQuotesMetadataToS3(tt.path, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadQuotesMetadataToS3() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 