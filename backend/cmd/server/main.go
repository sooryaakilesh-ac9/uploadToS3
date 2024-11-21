package main

import (
	"backend/internal/interface/http/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO define in environment file
	PORT := 8080

	mux := http.NewServeMux()

	router.RegisterHandlers(mux)

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), mux); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
