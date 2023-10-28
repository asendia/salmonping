package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/asendia/salmonping/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestGetPingAnomalies(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatal(err)
	}
	today09_59 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 59, 0, 0, loc)
	today10_01 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 1, 0, 0, loc)
	today10_11 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 11, 0, 0, loc)
	today10_21 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 21, 0, 0, loc)
	today10_31 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 31, 0, 0, loc)
	today10_41 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 41, 0, 0, loc)
	today19_51 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 51, 0, 0, loc)
	today19_59 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 59, 0, 0, loc)
	today20_01 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 1, 0, 0, loc)
	today20_11 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 11, 0, 0, loc)
	dayOfWeek := today10_01.Weekday()

	uuid1 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid2 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid3 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid4 := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	schedules := []db.SelectOnlineListingSchedulesRow{
		{OnlineListingID: uuid1, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_01), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_01), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
		{OnlineListingID: uuid2, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_01), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_01), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
		{OnlineListingID: uuid3, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_01), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_01), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
		{OnlineListingID: uuid4, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_01), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_01), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Grabfood: Resto A", Platform: "grabfood", Url: "https://grabfood.com/resto-a"},
	}
	// Define test cases
	testCases := []struct {
		name      string
		schedules []db.SelectOnlineListingSchedulesRow
		pings     []db.SelectOnlineListingPingsRow
		expected  []db.SelectOnlineListingPingsRow
	}{
		{
			name:      "anomalies: (1) Gofood: Resto A open -> closed, (2) Grabfood: Resto A closed first ping",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
		},
		{
			name:      "anomalies: Gofood: Resto A & C closed first ping",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
			},
		},
		{
			name:      "anomalies: Gofood: unknown 3 times in a row",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_41, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_41, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
			},
		},
		{
			name:      "no anomalies: open on first ping",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_59, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: closed 2x in a row",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: always open",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_31, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_21, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_11, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: closed outside operation shcedule",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: unkown outside operation shcedule",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_11, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_59, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_51, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getPingAnomalies(tc.schedules, tc.pings)
			if len(tc.expected) == 0 && len(actual) == 0 {
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				// Convert json to string
				actualJson, err := json.Marshal(actual)
				if err != nil {
					t.Fatal(err)
				}
				expectedJson, err := json.Marshal(tc.expected)
				if err != nil {
					t.Fatal(err)
				}
				t.Errorf("%s: expected %s\n\nBut got %s", tc.name, expectedJson, actualJson)
			}
		})
	}
}
