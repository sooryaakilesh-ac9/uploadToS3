package utils

import (
	"backend/pkg/images"
	"os"
	"path/filepath"
	"testing"
)

func TestImagesToJson(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	
	// Set all required environment variables
	os.Setenv("IMAGE_METADATA_FILENAME", "imagesMetadata.json")
	os.Setenv("IMAGE_METADATA_PATH", tmpDir)
	os.Setenv("IMAGE_METADATA_URL", "http://localhost:4566/test-bucket/images")
	os.Setenv("S3_REGION", "us-east-1")
	os.Setenv("S3_ENDPOINT", "http://localhost:4566")
	os.Setenv("S3_ID", "test")
	os.Setenv("S3_SECRET", "test")
	os.Setenv("S3_TOKEN", "")
	os.Setenv("S3_BUCKET_NAME", "test-bucket")
	os.Setenv("S3_IMAGES_DIR_PATH", "/images/")

	// Store original env values
	originalEnv := map[string]string{
		"IMAGE_METADATA_FILENAME": os.Getenv("IMAGE_METADATA_FILENAME"),
		"IMAGE_METADATA_PATH":     os.Getenv("IMAGE_METADATA_PATH"),
		"IMAGE_METADATA_URL":      os.Getenv("IMAGE_METADATA_URL"),
		"S3_REGION":              os.Getenv("S3_REGION"),
		"S3_ENDPOINT":            os.Getenv("S3_ENDPOINT"),
		"S3_ID":                  os.Getenv("S3_ID"),
		"S3_SECRET":              os.Getenv("S3_SECRET"),
		"S3_TOKEN":               os.Getenv("S3_TOKEN"),
		"S3_BUCKET_NAME":         os.Getenv("S3_BUCKET_NAME"),
		"S3_IMAGES_DIR_PATH":     os.Getenv("S3_IMAGES_DIR_PATH"),
	}

	// Restore original environment variables after test
	defer func() {
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	tests := []struct {
		name    string
		images  []images.Flyer
		wantErr bool
	}{
		{
			name:    "empty images list",
			images:  []images.Flyer{},
			wantErr: false,
		},
		{
			name: "valid images list",
			images: []images.Flyer{
				{
					Id: 1,
					Design: images.Design{
						TemplateId: "test1",
						Type:      "image",
						Tags:      []string{"test", "image"},
						Resolution: images.Resolution{
							Width:  1920,
							Height: 1080,
							Unit:   1,
						},
						FileFormat:  "JPEG",
						Orientation: "landscape",
					},
					Lang: "en-US",
					Url:  "test.jpg",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ImagesToJson(tt.images)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImagesToJson() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file was created
			expectedFile := filepath.Join(tmpDir, "imagesMetadata.json")
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) && !tt.wantErr {
				t.Errorf("ImagesToJson() failed to create file at %s", expectedFile)
			}
		})
	}
} 