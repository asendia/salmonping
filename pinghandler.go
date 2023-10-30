package main

import (
	"net/http"

	"github.com/asendia/salmonping/db"
	"github.com/gin-gonic/gin"
)

func pingHandler(c *gin.Context) {
	ctx := c.Request.Context()
	tx, conn, _, message, err := prepareDBConn(ctx)
	if conn != nil {
		defer conn.Close(ctx)
	}
	if tx != nil {
		// Commit everything
		defer tx.Commit(ctx)
	}
	if err != nil {
		log := map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": message,
		}
		logJson(log)
		c.JSON(http.StatusInternalServerError, gin.H(log))
		return
	}
	queries := db.New(tx)
	err = fetchListings(ctx, queries)
	if err != nil {
		log := map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error fetching listings",
		}
		logJson(log)
		c.JSON(http.StatusInternalServerError, gin.H(log))
		return
	}

	err = sendAlerts(ctx, queries)
	if err != nil {
		logJson(map[string]interface{}{
			"level":   "warning",
			"error":   err.Error(),
			"message": "Failed to send Telegram alert",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
