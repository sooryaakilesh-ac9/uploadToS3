package main

import (
	"backend/internal/interface/http/router"
	"backend/ops/db"
	"fmt"
	"log"
	"net/http"
)

// todo add the haredcoded values to the ENV file

func main() {
	PORT := 8080

	mux := http.NewServeMux()
	router.RegisterHandlers(mux)

	// test connection ping
	_, err := db.ConnectToDB()
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), mux); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
