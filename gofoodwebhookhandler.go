package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func gofoodWebhookHandler(c *gin.Context) {
	var payload GofoodWebhookPayload
	err := c.BindJSON(&payload)
	if err != nil {
		log := map[string]interface{}{
			"error":   err.Error(),
			"level":   "error",
			"message": "Error parsing body",
		}
		logJson(log)
		c.JSON(http.StatusBadRequest, gin.H(log))
		return
	}
	log := map[string]interface{}{
		"payload": payload,
		"headers": c.Request.Header,
		"level":   "info",
		"message": "Gofood Webhook received",
	}
	logJson(log)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
