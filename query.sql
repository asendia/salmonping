-- name: InsertOnlineListing :one
INSERT INTO online_listing (
  name,
  platform,
  url
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: InsertListing :one
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
) RETURNING *;

-- name: InsertPing :one
INSERT INTO salmon_ping (
  status,
  online_listing_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: SelectListings :many
SELECT
    ol.id,
    ol.created_at,
    ol.name,
    ol.platform,
    ol.url
FROM online_listing ol;

-- name: SelectOnlineListingSchedules :many
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
ORDER BY s.opening_time ASC;

-- name: SelectOnlineListingPings :many
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
    sp.created_at >= @start_date
    AND sp.created_at < @end_date
    AND sp.status = ANY(@statuses::text[])
ORDER BY sp.created_at DESC
LIMIT $1
OFFSET $2;
    