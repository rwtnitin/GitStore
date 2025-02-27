package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT string = "8080"

func main() {
	// Create a multiplexer
	var mux *http.ServeMux = http.NewServeMux()

	// Add routes to the multiplexer
	mux.HandleFunc("GET /health", HealthHandler)

	log.Println("Starting server on PORT", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux); err != nil {
		log.Fatalf("Error starting server on PORT %s: %s\n", PORT, err)
	}
}
