package main

import (
	"net/http"

	"github.com/asendia/salmonping/db"
	"github.com/gin-gonic/gin"
)

// pingHandler godoc
//
// @Summary		Ping and scrape online listings
// @Description	this endpoint is called by cloud scheduler
// @Tags		ping
// @Accept		json
// @Produce		json
// @Security	ApiKeyAuth
// @Success		200	{object}	DefaultResponse
// @Failure		401	{object}	DefaultErrorResponse
// @Failure		500	{object}	DefaultErrorResponse
// @Router		/ping [get]
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
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: message,
		}
		logJson(log.JSON())
		c.JSON(http.StatusInternalServerError, log)
		return
	}
	queries := db.New(tx)
	err = fetchListings(ctx, queries)
	if err != nil {
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error fetching listings",
		}
		logJson(log.JSON())
		c.JSON(http.StatusInternalServerError, log)
		return
	}

	err = sendAlerts(ctx, queries)
	if err != nil {
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error sending alerts",
		}
		logJson(log.JSON())
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "ok",
	})
}
