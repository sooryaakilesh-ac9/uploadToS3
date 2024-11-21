package router

import (
	"backend/internal/interface/http/handler"
	"backend/internal/interface/http/middleware"
	"net/http"
)

// register handlers
func RegisterHandlers(mux *http.ServeMux) {
	// handlers regarding quotes
	http.HandleFunc("/quotes/import", handler.HandleQuotesImport)
	mux.Handle("/quotes", middleware.QuotesJsonAndMethodValidator(
		http.HandlerFunc(handler.HandleQuotesUpload),
	))

	// handlers regarding images
	http.HandleFunc("/images/import", handler.HandleImagesImport)
	http.HandleFunc("/images", handler.HandleImagesUpload)
}
