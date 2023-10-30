package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")

		// Check if the provided API key matches the expected one.
		if key != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "API key is invalid",
			})
			return
		}

		c.Next()
	}
}
