// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Print("starting server...")
	http.HandleFunc("/api/ping", routePing)
	http.HandleFunc("/api/history", routeHistory)
	http.HandleFunc("/api/webhook/gofood", routeWebhookGofood)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
