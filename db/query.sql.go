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
	DayOfWeek       int32
	OpeningTime     pgtype.Time
	ClosingTime     pgtype.Time
	OnlineListingID pgtype.UUID
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
	Status          string
	OnlineListingID pgtype.UUID
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

const insertRestaurant = `-- name: InsertRestaurant :one
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

type InsertRestaurantParams struct {
	Name     string
	Platform string
	Url      string
}

func (q *Queries) InsertRestaurant(ctx context.Context, arg InsertRestaurantParams) (OnlineListing, error) {
	row := q.db.QueryRow(ctx, insertRestaurant, arg.Name, arg.Platform, arg.Url)
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

const selectOnlineListingSchedules = `-- name: SelectOnlineListingSchedules :many
SELECT
    ol.id,
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
    ol.id = s.restaurant_id
WHERE
    ol.id = $1
    AND s.day_of_week = $2
`

type SelectOnlineListingSchedulesParams struct {
	ID        pgtype.UUID
	DayOfWeek int32
}

type SelectOnlineListingSchedulesRow struct {
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	Name        string
	Platform    string
	Url         string
	DayOfWeek   int32
	OpeningTime pgtype.Time
	ClosingTime pgtype.Time
}

func (q *Queries) SelectOnlineListingSchedules(ctx context.Context, arg SelectOnlineListingSchedulesParams) ([]SelectOnlineListingSchedulesRow, error) {
	rows, err := q.db.Query(ctx, selectOnlineListingSchedules, arg.ID, arg.DayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectOnlineListingSchedulesRow
	for rows.Next() {
		var i SelectOnlineListingSchedulesRow
		if err := rows.Scan(
			&i.ID,
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
