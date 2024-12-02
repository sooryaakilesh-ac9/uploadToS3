package router

import (
	"backend/internal/interface/http/handler"
	"backend/internal/interface/http/middleware"
	"backend/ops/db"
	"log"
	"net/http"
)

// register handlers
func RegisterHandlers(mux *http.ServeMux) {
	// Initialize handlers
	database := db.NewPostgresDB()
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	imageHandler := handler.NewImageHandler(database)

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
	mux.Handle("/images/import", middleware.ImagesImport(
		http.HandlerFunc(imageHandler.HandleImagesImport),
	))
	// http method is checked
	// checks if the image is of valid type and is within the size limit
	mux.Handle("/images", middleware.ImageAndMethodValidator(
		http.HandlerFunc(imageHandler.HandleImagesUpload),
	))
}
