package router

import (
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"log"
	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, imageHandler *handler.ImageHandler) {
	// Create middleware chain
	chain := func(h http.Handler) http.Handler {
		return middleware.ErrorHandler(
			middleware.LogRequest(h),
		)
	}

	// Image routes with middleware chain
	mux.Handle("/images/import", chain(
		middleware.ImagesImport(
			http.HandlerFunc(imageHandler.HandleImagesImport),
		),
	))

	mux.Handle("/images", chain(
		middleware.ImageAndMethodValidator(
			http.HandlerFunc(imageHandler.HandleImagesUpload),
		),
	))

	log.Println("Routes registered successfully")
} 