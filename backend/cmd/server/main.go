package main

import (
	"backend/internal/interface/http/router"
	"backend/ops/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO define in environment file
	PORT := 8080

	mux := http.NewServeMux()

	router.RegisterHandlers(mux)

	_, err := db.ConnectToDB()
	if err != nil {
		log.Printf("%v", err)
	}

	// test fetching
	quote, err := db.FetchQuoteFromDB(1)
	if err != nil {
		log.Print(err)
	}

	log.Printf("%+v", quote)

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), mux); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
