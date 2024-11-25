package main

import (
	"backend/internal/interface/http/router"
	"backend/ops/db"
	"backend/utils"
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

	// test fetching
	quote, err := db.FetchQuoteFromDB(1)
	if err != nil {
		log.Print(err)
	}
	jsonQuote, err := utils.JsonHandler(quote)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(string(jsonQuote))

	log.Printf("Listening on PORT: %v...\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), mux); err != nil {
		log.Printf("Unable to start server on PORT: %v", PORT)
	}
}
