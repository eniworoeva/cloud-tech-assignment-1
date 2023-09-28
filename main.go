package main

import (
	"log"
	"net/http"
	"os"
	"osaigie/unisearcher/functions"
	"osaigie/unisearcher/handler"
	"github.com/rs/cors" // Import the CORS package
)

//Body taken from 02-JSON-demo
func main() {

	functions.GetUpTime()

	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler endpoints
	http.HandleFunc(handler.DEFAULT_PATH, handler.EmptyHandler)
	http.HandleFunc(handler.UNIINFO_PATH, handler.UniInfoHandler)
	http.HandleFunc(handler.NEIGHBOURUNIS_PATH, handler.NeighbourUnisHandler)
	http.HandleFunc(handler.DIAG_PATH, handler.DiagHandler)




	  // Create a CORS middleware with the desired CORS options
	  corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allow requests from any origin
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
    })

    // Wrap the default ServeMux with the CORS middleware
    http.Handle("/", corsHandler.Handler(http.DefaultServeMux))

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
