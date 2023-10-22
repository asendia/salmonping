package main

import (
	"context"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// pings should be ordered by created_at DESC
func getPingAnomalies(schedules []db.SelectOnlineListingSchedulesRow, pings []db.SelectOnlineListingPingsRow) []db.SelectOnlineListingPingsRow {
	var anomalies []db.SelectOnlineListingPingsRow
	// Get current and previous ping status for each online listing id/name
	currentListingPingMap := make(map[string]db.SelectOnlineListingPingsRow)
	previousListingPingMap := make(map[string]db.SelectOnlineListingPingsRow)
	for _, row := range pings {
		// Check if row.Name key does not exist in currentListingPingMap
		if _, ok := currentListingPingMap[row.Name]; !ok {
			currentListingPingMap[row.Name] = row
		} else if _, ok := previousListingPingMap[row.Name]; !ok {
			previousListingPingMap[row.Name] = row
		}
	}

	// Check if current ping status is different from previous ping status
	// If different, send message to Telegram
	for _, row := range schedules {
		// Check if row.Name key exists in currentListingPingMap
		current, isCurrentFound := currentListingPingMap[row.Name]
		if !isCurrentFound {
			continue
		}
		microsSinceMidnight := toMicrosSinceMidnight(current.CreatedAt.Time)

		// Check if current ping is within schedule
		if microsSinceMidnight < row.OpeningTime.Microseconds || microsSinceMidnight > row.ClosingTime.Microseconds {
			continue
		}
		// Check if row.Name key exists in previousListingPingMap
		previous, isPreviousFound := previousListingPingMap[row.Name]
		// Check if first ping within schedule or current ping status is different from previous ping status
		if isPreviousFound && current.Status == previous.Status {
			continue
		}
		if current.Status == "open" {
			continue
		}
		anomalies = append(anomalies, current)
	}
	return anomalies
}

func getTodaySchedules(ctx context.Context, queries *db.Queries) ([]db.SelectOnlineListingSchedulesRow, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}
	dayOfWeek := time.Now().In(loc).Weekday()
	schedules, err := queries.SelectOnlineListingSchedules(ctx, int32(dayOfWeek))
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func getTodayPings(ctx context.Context, queries *db.Queries) ([]db.SelectOnlineListingPingsRow, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}
	today := time.Now().In(loc)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)
	listingPings, err := queries.SelectOnlineListingPings(ctx, db.SelectOnlineListingPingsParams{
		Limit:     100,
		Offset:    0,
		StartDate: pgtype.Timestamptz{Time: today, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: time.Now().In(loc), Valid: true},
	})
	return listingPings, err
}
