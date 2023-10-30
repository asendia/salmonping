package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// /api/history?page=1&start=2021-01-01&end=2021-01-31&status=closed,unknown&platform=gofood,grabfood&name=Kebon%20Jeruk,Sudirman
func historyHandler(c *gin.Context) {
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
		log := map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": message,
		}
		logJson(log)
		c.JSON(http.StatusInternalServerError, gin.H(log))
		return
	}

	var payload HistoryPayload
	if err := c.ShouldBindQuery(&payload); err != nil {
		log := map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error binding payload",
			"query":   c.Request.URL.RawQuery,
		}
		logJson(log)
		c.JSON(http.StatusBadRequest, gin.H(log))
		return
	}

	// Query string params
	queries := db.New(tx)
	limit := int32(100)
	if payload.Page == 0 {
		payload.Page = 1
	}
	offset := (int32(payload.Page) - 1) * limit
	startDate := parseJakartaDate(payload.StartDate, time.Now().Add(-24*7*time.Hour))
	endDate := parseJakartaDate(payload.EndDate, time.Now())
	// Add 1 day to endDate
	endDate = endDate.Add(24 * time.Hour)
	pgStartDate := pgtype.Timestamptz{Time: startDate, Valid: true}
	pgEndDate := pgtype.Timestamptz{Time: endDate, Valid: true}
	if payload.Name == "" {
		payload.Name = "Haji Nawi,Kebon Jeruk,Sudirman"
	}
	if payload.Platform == "" {
		payload.Platform = "gofood,grabfood"
	}
	if payload.Status == "" {
		payload.Status = "open,closed,unknown"
	}

	listingPings, err := queries.SelectOnlineListingPings(ctx, db.SelectOnlineListingPingsParams{
		EndDate:   pgEndDate,
		Limit:     limit,
		Names:     strings.Split(payload.Name, ","),
		Offset:    offset,
		Platforms: strings.Split(payload.Platform, ","),
		StartDate: pgStartDate,
		Statuses:  strings.Split(payload.Status, ","),
	})
	if err != nil {
		log := map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error selecting listing pings",
		}
		logJson(log)
		c.JSON(http.StatusInternalServerError, gin.H(log))
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"listing_pings": listingPings,
	})
}

type HistoryPayload struct {
	EndDate   string `form:"end"`
	Name      string `form:"name"`
	Page      int    `form:"page"`
	Platform  string `form:"platform"`
	StartDate string `form:"start"`
	Status    string `form:"status"`
}
