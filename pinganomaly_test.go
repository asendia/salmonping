package main

import (
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
	today09_55 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 55, 0, 0, loc)
	today10_00 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 0, 0, 0, loc)
	today10_05 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 5, 0, 0, loc)
	today10_30 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 30, 0, 0, loc)
	today10_55 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 55, 0, 0, loc)
	today19_30 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 30, 0, 0, loc)
	today19_55 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 55, 0, 0, loc)
	today20_00 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 0, 0, 0, loc)
	today20_05 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 5, 0, 0, loc)
	dayOfWeek := today10_00.Weekday()

	uuid1 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid2 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid3 := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	uuid4 := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	schedules := []db.SelectOnlineListingSchedulesRow{
		{OnlineListingID: uuid1, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_00), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_00), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
		{OnlineListingID: uuid2, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_00), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_00), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
		{OnlineListingID: uuid3, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_00), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_00), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
		{OnlineListingID: uuid4, OpeningTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today10_00), Valid: true}, ClosingTime: pgtype.Time{Microseconds: toMicrosSinceMidnight(today20_00), Valid: true}, DayOfWeek: int32(dayOfWeek), Name: "Grabfood: Resto A", Platform: "grabfood", Url: "https://grabfood.com/resto-a"},
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
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
		},
		{
			name:      "anomalies: Gofood: Resto A & C closed first ping",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
			},
		},
		{
			name:      "anomalies: Gofood: Resto C unknown 3 times in a row",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
			},
			expected: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-c"},
			},
		},
		{
			name:      "no anomalies: open on first ping",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today09_55, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: closed 2x in a row",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "closed", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: always open",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_55, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_30, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today10_05, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: closed outside operation shcedule",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "closed", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "closed", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "closed", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "open", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "open", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "open", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
			},
			expected: []db.SelectOnlineListingPingsRow{},
		},
		{
			name:      "no anomalies: unkown outside operation shcedule",
			schedules: schedules,
			pings: []db.SelectOnlineListingPingsRow{
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today20_05, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_55, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "open", Name: "Gofood: Resto A", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "unknown", Name: "Gofood: Resto B", Platform: "gofood", Url: "https://gofood.com/resto-b"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "unknown", Name: "Gofood: Resto C", Platform: "gofood", Url: "https://gofood.com/resto-a"},
				{OnlineListingID: uuid1, CreatedAt: pgtype.Timestamptz{Time: today19_30, Valid: true}, Status: "unknown", Name: "Grabfood: Resto A", Platform: "gofood", Url: "https://grabfood.com/resto-a"},
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
				t.Errorf("%s: expected %v, but got %v", tc.name, tc.expected, actual)
			}
		})
	}
}
