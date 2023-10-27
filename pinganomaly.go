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
	pingMap := make(map[string][]db.SelectOnlineListingPingsRow)
	for _, row := range pings {
		// Check if row.Name key does not exist in currentListingPingMap
		if _, ok := pingMap[row.Name]; !ok {
			pingMap[row.Name] = []db.SelectOnlineListingPingsRow{row}
		} else {
			pingMap[row.Name] = append(pingMap[row.Name], row)
		}
	}

	// Check if current ping status is different from previous ping status
	// If different, send message to Telegram
	for _, row := range schedules {
		// Check if row.Name key exists in currentListingPingMap
		p, ok := pingMap[row.Name]
		if !ok || p[0].Status == "open" {
			continue
		}

		if isBetweenTime(p[0].CreatedAt.Time, row.OpeningTime, row.ClosingTime) {
			continue
		}

		isFirstClosed := p[0].Status == "closed" && (len(p) == 1 || (len(p) > 1 && !isBetweenTime(p[1].CreatedAt.Time, row.OpeningTime, row.ClosingTime)))
		isSwitchingToClosed := p[0].Status == "closed" && len(p) > 1 && p[1].Status != "closed"
		if isFirstClosed || isSwitchingToClosed {
			anomalies = append(anomalies, p[0])
			continue
		}

		unknownCombo := 0
		for i := 0; i < len(p); i++ {
			if p[i].Status == "unknown" {
				unknownCombo++
			} else {
				break
			}
		}
		if unknownCombo == 3 {
			anomalies = append(anomalies, p[0])
		}
	}
	return anomalies
}

func isBetweenTime(t time.Time, start pgtype.Time, end pgtype.Time) bool {
	microsSinceMidnight := toMicrosSinceMidnight(t)
	return microsSinceMidnight >= start.Microseconds && microsSinceMidnight <= end.Microseconds
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
	tomorrow := today.AddDate(0, 0, 1)
	listingPings, err := queries.SelectOnlineListingPings(ctx, db.SelectOnlineListingPingsParams{
		Limit:     100,
		Offset:    0,
		StartDate: pgtype.Timestamptz{Time: today, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: tomorrow, Valid: true},
		Statuses:  []string{"open", "closed", "unknown"},
	})
	return listingPings, err
}
