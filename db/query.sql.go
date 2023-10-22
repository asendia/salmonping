// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertListing = `-- name: InsertListing :one
INSERT INTO schedule (
  day_of_week,
  opening_time,
  closing_time,
  online_listing_id
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING id, online_listing_id, day_of_week, opening_time, closing_time
`

type InsertListingParams struct {
	DayOfWeek       int32       `json:"day_of_week"`
	OpeningTime     pgtype.Time `json:"opening_time"`
	ClosingTime     pgtype.Time `json:"closing_time"`
	OnlineListingID pgtype.UUID `json:"online_listing_id"`
}

func (q *Queries) InsertListing(ctx context.Context, arg InsertListingParams) (Schedule, error) {
	row := q.db.QueryRow(ctx, insertListing,
		arg.DayOfWeek,
		arg.OpeningTime,
		arg.ClosingTime,
		arg.OnlineListingID,
	)
	var i Schedule
	err := row.Scan(
		&i.ID,
		&i.OnlineListingID,
		&i.DayOfWeek,
		&i.OpeningTime,
		&i.ClosingTime,
	)
	return i, err
}

const insertOnlineListing = `-- name: InsertOnlineListing :one
INSERT INTO online_listing (
  name,
  platform,
  url
) VALUES (
  $1,
  $2,
  $3
) RETURNING id, created_at, name, platform, url
`

type InsertOnlineListingParams struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

func (q *Queries) InsertOnlineListing(ctx context.Context, arg InsertOnlineListingParams) (OnlineListing, error) {
	row := q.db.QueryRow(ctx, insertOnlineListing, arg.Name, arg.Platform, arg.Url)
	var i OnlineListing
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Platform,
		&i.Url,
	)
	return i, err
}

const insertPing = `-- name: InsertPing :one
INSERT INTO salmon_ping (
  status,
  online_listing_id
) VALUES (
  $1,
  $2
) RETURNING id, created_at, status, online_listing_id
`

type InsertPingParams struct {
	Status          string      `json:"status"`
	OnlineListingID pgtype.UUID `json:"online_listing_id"`
}

func (q *Queries) InsertPing(ctx context.Context, arg InsertPingParams) (SalmonPing, error) {
	row := q.db.QueryRow(ctx, insertPing, arg.Status, arg.OnlineListingID)
	var i SalmonPing
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Status,
		&i.OnlineListingID,
	)
	return i, err
}

const selectListings = `-- name: SelectListings :many
SELECT
    ol.id,
    ol.created_at,
    ol.name,
    ol.platform,
    ol.url
FROM online_listing ol
`

func (q *Queries) SelectListings(ctx context.Context) ([]OnlineListing, error) {
	rows, err := q.db.Query(ctx, selectListings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OnlineListing
	for rows.Next() {
		var i OnlineListing
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Platform,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOnlineListingPings = `-- name: SelectOnlineListingPings :many
SELECT
    sp.id AS salmon_ping_id,
    sp.created_at,
    sp.status,
    ol.id AS online_listing_id,
    ol.name,
    ol.platform,
    ol.url
FROM salmon_ping sp
JOIN online_listing ol
ON
    sp.online_listing_id = ol.id
WHERE
    sp.created_at >= $3
    AND sp.created_at <= $4
    AND sp.status = ANY($5::text[])
ORDER BY sp.created_at DESC
LIMIT $1
OFFSET $2
`

type SelectOnlineListingPingsParams struct {
	Limit     int32              `json:"limit"`
	Offset    int32              `json:"offset"`
	StartDate pgtype.Timestamptz `json:"start_date"`
	EndDate   pgtype.Timestamptz `json:"end_date"`
	Statuses  []string           `json:"statuses"`
}

type SelectOnlineListingPingsRow struct {
	SalmonPingID    pgtype.UUID        `json:"salmon_ping_id"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	Status          string             `json:"status"`
	OnlineListingID pgtype.UUID        `json:"online_listing_id"`
	Name            string             `json:"name"`
	Platform        string             `json:"platform"`
	Url             string             `json:"url"`
}

func (q *Queries) SelectOnlineListingPings(ctx context.Context, arg SelectOnlineListingPingsParams) ([]SelectOnlineListingPingsRow, error) {
	rows, err := q.db.Query(ctx, selectOnlineListingPings,
		arg.Limit,
		arg.Offset,
		arg.StartDate,
		arg.EndDate,
		arg.Statuses,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectOnlineListingPingsRow
	for rows.Next() {
		var i SelectOnlineListingPingsRow
		if err := rows.Scan(
			&i.SalmonPingID,
			&i.CreatedAt,
			&i.Status,
			&i.OnlineListingID,
			&i.Name,
			&i.Platform,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOnlineListingSchedules = `-- name: SelectOnlineListingSchedules :many
SELECT
    ol.id AS online_listing_id,
    ol.created_at,
    ol.name,
    ol.platform,
    ol.url,
    s.day_of_week,
    s.opening_time,
    s.closing_time
FROM online_listing ol
JOIN schedule s
ON
    ol.id = s.online_listing_id
WHERE
    s.day_of_week = $1
ORDER BY s.opening_time ASC
`

type SelectOnlineListingSchedulesRow struct {
	OnlineListingID pgtype.UUID        `json:"online_listing_id"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	Name            string             `json:"name"`
	Platform        string             `json:"platform"`
	Url             string             `json:"url"`
	DayOfWeek       int32              `json:"day_of_week"`
	OpeningTime     pgtype.Time        `json:"opening_time"`
	ClosingTime     pgtype.Time        `json:"closing_time"`
}

func (q *Queries) SelectOnlineListingSchedules(ctx context.Context, dayOfWeek int32) ([]SelectOnlineListingSchedulesRow, error) {
	rows, err := q.db.Query(ctx, selectOnlineListingSchedules, dayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectOnlineListingSchedulesRow
	for rows.Next() {
		var i SelectOnlineListingSchedulesRow
		if err := rows.Scan(
			&i.OnlineListingID,
			&i.CreatedAt,
			&i.Name,
			&i.Platform,
			&i.Url,
			&i.DayOfWeek,
			&i.OpeningTime,
			&i.ClosingTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
