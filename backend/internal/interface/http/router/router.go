package router

import (
	"backend/internal/interface/http/handler"
	"net/http"
)

// register handlers
func RegisterHandlers() {
	http.HandleFunc("/quotes/import", handler.HandleQuotesImport)
	http.HandleFunc("/quotes", handler.HandleQuotesUpload)
}
