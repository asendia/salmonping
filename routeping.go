package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/asendia/salmonping/db"
)

func routePing(w http.ResponseWriter, r *http.Request) {
	// Check API key
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != os.Getenv("API_KEY") {
		j, _ := logJson(map[string]interface{}{
			"level":   "warning",
			"message": "Invalid API key",
		})
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
		return
	}

	ctx := r.Context()
	tx, conn, _, message, err := prepareDBConn(ctx)
	if conn != nil {
		defer conn.Close(ctx)
	}
	if tx != nil {
		// Commit everything
		defer tx.Commit(ctx)
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

	err = fetchListings(ctx, queries)
	if err != nil {
		j, _ := logJson(map[string]interface{}{
			"level":   "error",
			"error":   err.Error(),
			"message": "Error fetching listings",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
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

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ok"}`)
}
