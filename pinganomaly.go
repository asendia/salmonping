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
		if _, ok := pingMap[row.Url]; !ok {
			pingMap[row.Url] = []db.SelectOnlineListingPingsRow{row}
		} else {
			pingMap[row.Url] = append(pingMap[row.Url], row)
		}
	}

	for _, row := range schedules {
		// Check if row.Name key exists in currentListingPingMap
		p, ok := pingMap[row.Url]
		if !ok || p[0].Status == "open" || !isBetweenTime(p[0].CreatedAt.Time, row.OpeningTime, row.ClosingTime) {
			continue
		}

		// Alert if within operational hours, status is changed from open to closed
		// while ignoring unknown statuses during operational hours
		isClosedAnomaly := false
		if p[0].Status == "closed" {
			isClosedAnomaly = true
			for i := 1; i < len(p); i++ {
				if !isBetweenTime(p[i].CreatedAt.Time, row.OpeningTime, row.ClosingTime) {
					break
				}
				if p[i].Status == "closed" {
					isClosedAnomaly = false
					break
				} else if p[i].Status == "open" {
					break
				}
			}
		}
		if isClosedAnomaly {
			anomalies = append(anomalies, p[0])
			continue
		}

		// Alert if unknown 3 times in a row during operational hours
		unknownCombo := 0
		for i := 0; i < len(p); i++ {
			if p[i].Status == "unknown" && isBetweenTime(p[i].CreatedAt.Time, row.OpeningTime, row.ClosingTime) {
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
		EndDate:   pgtype.Timestamptz{Time: tomorrow, Valid: true},
		Limit:     100,
		Names:     []string{"Haji Nawi", "Kebon Jeruk", "Sudirman", "Tanjung Duren"},
		Offset:    0,
		Platforms: []string{"gofood", "grabfood"},
		StartDate: pgtype.Timestamptz{Time: today, Valid: true},
		Statuses:  []string{"open", "closed", "unknown"},
	})
	return listingPings, err
}
