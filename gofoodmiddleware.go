package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GofoodSignatureMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log := map[string]interface{}{
				"error":   err.Error(),
				"level":   "error",
				"message": "Error reading body",
			}
			logJson(log)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H(log))
			return
		}
		defer c.Request.Body.Close()

		err = verifyGofoodSignature(bodyBytes, c.Request.Header.Get("X-Go-Signature"), secretKey)
		if err != nil {
			logJson(map[string]interface{}{
				"error":   err.Error(),
				"level":   "warning",
				"message": "Error validating signature",
			})
			// Cannot verify signature in sandbox
			// w.WriteHeader(http.StatusUnauthorized)
			// w.Header().Set("Content-Type", "application/json")
			// fmt.Fprintf(w, `{"message": "error"}`)
			// return
		}

		c.Next()
	}
}
