package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// /api/history?page=1&start=2021-01-01&end=2021-01-31&status=closed,unknown&platform=gofood,grabfood&name=Kebon%20Jeruk,Sudirman
func historyHandler(w http.ResponseWriter, r *http.Request) {
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

	// Query string params
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
		startDate = time.Now().In(loc)
		// Reset to 00:00:00
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		// Minues 7 days
		startDate = startDate.AddDate(0, 0, -7)
	}
	// Parse date from query string "end" if available
	endStr := r.URL.Query().Get("end")
	var endDate time.Time
	if endStr != "" {
		endDate, err = time.ParseInLocation("2006-01-02", endStr, loc)
	}
	if endStr == "" || err != nil {
		endDate = time.Now().In(loc)
		// Reset to 00:00:00
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)
	}
	// Add 1 day to endDate
	endDate = endDate.Add(24 * time.Hour)
	status := r.URL.Query().Get("status")
	// Split status by ","
	if status == "" {
		status = "open,closed,unknown"
	}
	statuses := strings.Split(status, ",")
	platform := r.URL.Query().Get("platform")
	// Split platform by ","
	if platform == "" {
		platform = "gofood,grabfood"
	}
	platforms := strings.Split(platform, ",")
	name := r.URL.Query().Get("name")
	// Split name by ","
	if name == "" {
		name = "Haji Nawi,Kebon Jeruk,Sudirman"
	}
	names := strings.Split(name, ",")

	pgStartDate := pgtype.Timestamptz{Time: startDate, Valid: true}
	pgEndDate := pgtype.Timestamptz{Time: endDate, Valid: true}
	listingPings, err := queries.SelectOnlineListingPings(ctx, db.SelectOnlineListingPingsParams{
		EndDate:   pgEndDate,
		Limit:     limit,
		Names:     names,
		Offset:    offset,
		Platforms: platforms,
		StartDate: pgStartDate,
		Statuses:  statuses,
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
