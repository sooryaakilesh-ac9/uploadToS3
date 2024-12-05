package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandleImagesUpload(t *testing.T) {
	// Create a temporary test image
	tmpImage := "test.jpg"
	if err := os.WriteFile(tmpImage, []byte("fake image data"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpImage)

	tests := []struct {
		name           string
		setupRequest   func() (*http.Request, error)
		expectedStatus int
	}{
		{
			name: "valid image upload",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("image", filepath.Base(tmpImage))
				if err != nil {
					return nil, err
				}
				imageData, err := os.ReadFile(tmpImage)
				if err != nil {
					return nil, err
				}
				part.Write(imageData)
				writer.Close()

				req := httptest.NewRequest("POST", "/images", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing image file",
			setupRequest: func() (*http.Request, error) {
				return httptest.NewRequest("POST", "/images", nil), nil
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.setupRequest()
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := NewImageHandler(nil) // You might need to mock the database
			handler.HandleImagesUpload(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
} 