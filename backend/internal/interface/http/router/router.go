package router

import (
	"backend/internal/interface/http/handler"
	"net/http"
)

// register handlers
func RegisterHandlers() {
	// handlers regarding quotes
	http.HandleFunc("/quotes/import", handler.HandleQuotesImport)
	http.HandleFunc("/quotes", handler.HandleQuotesUpload)

	// handlers regarding images
	http.HandleFunc("/images/import", handler.HandleImagesImport)
	http.HandleFunc("/images", handler.HandleImagesUpload)
}
