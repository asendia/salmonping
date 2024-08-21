package main

import (
	"net/http"
	"strings"

	"github.com/asendia/salmonping/db"
	"github.com/gin-gonic/gin"
)

// storesHandler godoc
//
// @Summary		Get list of stores
// @Description	get list of stores based on query string params
// @Tags		ping
// @Accept		json
// @Produce		json
// @Param		enable_ping	query		string	false	"Enable ping, true|false"
// @Param		name		query		string	false	"Names (comma spearated)"
// @Param		platform	query		string	false	"Platforms (comma spearated)"
// @Param		status		query		string	false	"Statuses (comma spearated)"
// @Success		200			{object}	StoresResponse
// @Failure		400			{object}	DefaultErrorResponse
// @Failure		500			{object}	DefaultErrorResponse
// @Router		/stores	[get]
func storesHandler(c *gin.Context) {
	// Prepare db connection
	ctx := c.Request.Context()
	tx, conn, _, message, err := prepareDBConn(ctx)
	if conn != nil {
		defer conn.Close(ctx)
	}
	if tx != nil {
		// Rollback everything
		defer tx.Rollback(ctx)
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

	var payload StoresPayload
	if err := c.ShouldBindQuery(&payload); err != nil {
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error binding payload",
			Query:   c.Request.URL.RawQuery,
		}
		logJson(log.JSON())
		c.JSON(http.StatusBadRequest, log)
		return
	}

	// Query string params
	queries := db.New(tx)
	var names = filterEmptyStrings(strings.Split(payload.Name, ","))
	var platforms = filterEmptyStrings(strings.Split(payload.Platform, ","))
	var statuses = filterEmptyStrings(strings.Split(payload.Status, ","))

	stores, err := queries.SelectListings(ctx, db.SelectListingsParams{
		EnablePing: []bool{true, false},
		Names:      names,
		Platforms:  platforms,
		Statuses:   statuses,
	})
	if err != nil {
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error selecting listing pings",
		}
		logJson(log.JSON())
		c.JSON(http.StatusInternalServerError, log)
		return
	}
	c.JSON(http.StatusOK, StoresResponse{
		Stores: stores,
	})
}

type StoresPayload struct {
	EnablePing string `form:"enable_ping"`
	Name       string `form:"name"`
	Platform   string `form:"platform"`
	Status     string `form:"status"`
	URL        string `form:"url"`
}

type StoresResponse struct {
	Stores []db.SelectListingsRow `json:"stores"`
}
