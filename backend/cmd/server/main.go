package main

import (
	"log"
	"net/http"

	"OurChat/internal/api"
)

func main() {
	// Create a new API server
	server := api.NewServer()

	// Configure the server routes
	server.SetupRoutes()

	// Start the server
	log.Println("Starting OurChat server on :8080")
	if err := http.ListenAndServe(":8080", server.Router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
