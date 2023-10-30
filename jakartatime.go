package main

import "time"

func parseJakartaDate(val string, fallback time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return fallback
	}
	t, err := time.ParseInLocation("2006-01-02", val, loc)
	if err != nil {
		return fallback
	}
	return t
}
