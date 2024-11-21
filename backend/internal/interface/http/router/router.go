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
	// checks if all the fields are present and if the JSON is valid
	// also http method is checked
	mux.Handle("/quotes", middleware.QuotesJsonAndMethodValidator(
		http.HandlerFunc(handler.HandleQuotesUpload),
	))

	// handlers regarding images
	http.HandleFunc("/images/import", handler.HandleImagesImport)
	mux.Handle("/images", middleware.ImageAndMethodValidator(
		http.HandlerFunc(handler.HandleImagesUpload),
	))
}
