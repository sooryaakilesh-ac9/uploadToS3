package middleware

import (
	"net/http"
	"os"
	"strconv"
)

// checks if the image is of valid type
// allowed image types (jpg, png)
// http methods are checked
func ImageAndMethodValidator(next http.Handler) http.Handler {
	// checks if the HTTP method is valid or not
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// setup ENV file to contain these details
		// Maximum upload size in MB
		maxUploadMB := 10 // Default value
		if envMB := os.Getenv("MAX_UPLOAD_MB"); envMB != "" {
			if val, err := strconv.ParseInt(envMB, 10, 64); err == nil {
				maxUploadMB = int(val)
			}
		}

		// Convert MB to bytes
		maxUploadSize := int64(maxUploadMB << 20)

		// Validate the methods (only POST and GET are allowed)
		if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
			http.Error(w, "Method Not Allowed. Only POST and GET are allowed.", http.StatusMethodNotAllowed)
			return
		}

		// Limits the size of the image
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			http.Error(w, "File too large. Maximum allowed size is 10MB.", http.StatusRequestEntityTooLarge)
			return
		}

		// Retrieve the file from the form
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Invalid file. Please upload a valid image.", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read the first 512 bytes to check the file type
		buffer := make([]byte, 512) // 512 bytes are enough to detect the file type
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Unable to read file. Please try again.", http.StatusInternalServerError)
			return
		}

		// Detect the file type based on the first 512 bytes
		fileType := http.DetectContentType(buffer)
		if fileType != "image/jpeg" && fileType != "image/png" {
			http.Error(w, "Invalid file type. Only JPG and PNG files are allowed.", http.StatusBadRequest)
			return
		}

		// If everything is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}
