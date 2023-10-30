// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Print("starting server...")
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://salmonfit.com", "https://salmonfit.id", "http://localhost:5173"}
	r.GET("/api/history", cors.New(corsConfig), historyHandler)

	apiKeyAuthMiddleware := APIKeyAuthMiddleware(os.Getenv("API_KEY"))
	r.GET("/api/ping", apiKeyAuthMiddleware, pingHandler)

	gofoodSignatureMiddleware := GofoodSignatureMiddleware(os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY"))
	r.POST("/api/webhook/gofood", gofoodSignatureMiddleware, gofoodWebhookHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
