package middleware

import (
	"net/http"
	"os"
	"strconv"
)

func ImageAndMethodValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Method == http.MethodPost {
			maxSize := getMaxUploadSize()
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			if err := r.ParseMultipartForm(maxSize); err != nil {
				http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func ImagesImport(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getMaxUploadSize() int64 {
	maxMB := 10 // Default 10MB
	if size := os.Getenv("MAX_UPLOAD_MB"); size != "" {
		if parsed, err := strconv.Atoi(size); err == nil {
			maxMB = parsed
		}
	}
	return int64(maxMB) << 20 // Convert to bytes
} 