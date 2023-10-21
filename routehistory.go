package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// /api/history?page=1&start=2021-01-01&end=2021-01-31
func routeHistory(w http.ResponseWriter, r *http.Request) {
	// Only allow API calls from salmonfit.com & salmonfit.id
	origin := r.Header.Get("Origin")
	if origin != "https://salmonfit.com" && origin != "https://salmonfit.id" && origin != "http://localhost:5173" {
		j, _ := logJson(map[string]interface{}{
			"level":   "warning",
			"message": "Invalid origin",
		})
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}

	// Write CORS headers
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD, POST")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Prepare db connection
	ctx := r.Context()
	tx, conn, _, message, err := prepareDBConn(ctx)
	if conn != nil {
		defer conn.Close(ctx)
	}
	if tx != nil {
		// Rollback everything
		defer tx.Rollback(ctx)
	}
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": message,
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	queries := db.New(tx)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit := int32(100)
	offset := (int32(page) - 1) * limit
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error loading timezone",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	// Parse date from query string "start" if available
	startStr := r.URL.Query().Get("start")
	var startDate time.Time
	if startStr != "" {
		startDate, err = time.ParseInLocation("2006-01-02", startStr, loc)
	}
	if startStr == "" || err != nil {
		startDate = time.Now().Add(-24 * time.Hour)
		// Reset to 00:00:00
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	}
	// Parse date from query string "end" if available
	endStr := r.URL.Query().Get("end")
	var endDate time.Time
	if endStr != "" {
		endDate, err = time.ParseInLocation("2006-01-02", endStr, loc)
	}
	if endStr == "" || err != nil {
		endDate = time.Now()
		// Reset to 23:59:59
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, loc)
	}

	var pgStartDate, pgEndDate pgtype.Timestamptz
	pgStartDate.Time = startDate
	pgStartDate.Valid = true
	pgEndDate.Time = endDate
	pgEndDate.Valid = true
	listingPings, err := queries.SelectOnlineListingPings(ctx, db.SelectOnlineListingPingsParams{
		EndDate:   pgEndDate,
		Limit:     limit,
		Offset:    offset,
		StartDate: pgStartDate,
	})
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error selecting listing pings",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	// return listingPings as JSON
	j, err := json.Marshal(map[string]interface{}{
		"listing_pings": listingPings,
	})
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error marshaling JSON",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(j))
}
