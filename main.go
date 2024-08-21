// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"log"
	"os"

	"github.com/asendia/salmonping/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						Salmon Ping API
// @version						1.0
// @description					Online listing status checker by Salmon Fit.
// @contact.name				Salmon Ping
// @contact.url					https://salmonfit.com
// @license.name				Apache 2.0
// @license.url					http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath					/api
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-Key
// @description					Static API key for authentication
// @securityDefinitions.apikey	GofoodSignature
// @in							header
// @name						X-Go-Signature
// @description					HMAC sha256 signature based on content body and secret key
// @externalDocs.description	OpenAPI
// @externalDocs.url			/swagger/index.html
func main() {
	godotenv.Load()
	// Init swagger docs, it needs to be imported in main.go
	docs.SwaggerInfo.Title = "Salmon Ping API"

	log.Print("starting server...")
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://salmonfit.com", "https://salmonfit.id", "http://localhost:5173"}
	r.GET("/api/history", cors.New(corsConfig), historyHandler)
	r.GET("/api/stores", cors.New(corsConfig), storesHandler)

	apiKeyAuthMiddleware := APIKeyAuthMiddleware(os.Getenv("API_KEY"))
	r.GET("/api/ping", apiKeyAuthMiddleware, pingHandler)

	gofoodSignatureMiddleware := GofoodSignatureMiddleware(os.Getenv("GOFOOD_NOTIFICATION_SECRET_KEY"))
	r.POST("/api/webhook/gofood", gofoodSignatureMiddleware, gofoodWebhookHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
