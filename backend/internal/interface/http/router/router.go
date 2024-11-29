package router

import (
	"backend/internal/interface/http/handler"
	"backend/internal/interface/http/middleware"
	"net/http"
)

// register handlers
func RegisterHandlers(mux *http.ServeMux) {
	// handlers regarding quotes
	// http.HandleFunc("/quotes/import", handler.HandleQuotesImport)
	mux.Handle("/quotes/import", middleware.CheckQuotesLink(
		http.HandlerFunc(handler.HandleQuotesImport)),
	)
	// checks if all the fields are present and if the JSON is valid
	// also http method is checked
	mux.Handle("/quotes", middleware.QuotesJsonAndMethodValidator(
		http.HandlerFunc(handler.HandleQuotesUpload),
	))

	// handlers regarding images

	// pending (google drive or download and import?)
	http.HandleFunc("/images/import", handler.HandleImagesImport)
	mux.Handle("/images/import", middleware.ImagesImport(
		http.HandlerFunc(handler.HandleImagesImport),
	))
	// http method is checked
	// checks if the image is of valid type and is within the size limit
	mux.Handle("/images", middleware.ImageAndMethodValidator(
		http.HandlerFunc(handler.HandleImagesUpload),
	))
}
