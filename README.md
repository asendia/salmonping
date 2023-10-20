# Salmon Ping
Online listing status checker by Salmon Fit

## Prerequisites
- [Go](go) 1.7+, but I use 1.21.3
- [Sqlx](https://docs.sqlc.dev/en/latest/overview/install.html) for development

## Setup
```sh
cp .env-template .env
# Then update the variable accordingly
```

## Run
```sh
go run .
# Run at http://localhost:8080/api/ping
```

## Sqlc
1. Modify schema.sql or query.sql
2. Run `sqlc generate`
3. Files in `db` dir should be updated
