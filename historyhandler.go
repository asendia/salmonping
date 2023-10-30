package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// historyHandler godoc
//
// @Summary		Show salmon ping history
// @Description	get ping history based on query string params
// @Tags		ping
// @Accept		json
// @Produce		json
// @Param		page		query		int		false	"Page"							default(1)
// @Param		start		query		string	false	"Start Date (inclusive)"		example("2023-10-28")
// @Param		end			query		string	false	"End Date (inclusive)"			example("2023-10-31")
// @Param		name		query		string	false	"Names (comma spearated)"		default(Haji Nawi,Kebon Jeruk,Sudirman)
// @Param		platform	query		string	false	"Platforms (comma spearated)"	default(gofood,grabfood)
// @Param		status		query		string	false	"Statuses (comma spearated)"	default(open,closed,unknown)
// @Success		200			{object}	HistoryResponse
// @Failure		400			{object}	DefaultErrorResponse
// @Failure		500			{object}	DefaultErrorResponse
// @Router		/history	[get]
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
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: message,
		}
		logJson(log.JSON())
		c.JSON(http.StatusInternalServerError, log)
		return
	}

	var payload HistoryPayload
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
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error selecting listing pings",
		}
		logJson(log.JSON())
		c.JSON(http.StatusInternalServerError, log)
		return
	}
	c.JSON(http.StatusOK, HistoryResponse{
		ListingPings: listingPings,
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

type HistoryResponse struct {
	ListingPings []db.SelectOnlineListingPingsRow `json:"listing_pings"`
}
