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

	router.RegisterHandlers()

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), nil); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
