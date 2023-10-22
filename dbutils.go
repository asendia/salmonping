package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func prepareDBConn(ctx context.Context) (tx pgx.Tx, conn *pgx.Conn, config *pgx.ConnConfig, message string, err error) {
	config, err = pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		message = "Error parsing database URL"
		return
	}
	// Since neondb doesn't support prepared statements, we need to use simple protocol
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Connect to DB
	conn, err = pgx.ConnectConfig(ctx, config)
	if err != nil {
		message = "Error connecting to database"
		return
	}

	// Start transaction
	tx, err = conn.Begin(ctx)
	if err != nil {
		message = "Error starting transaction"
		return
	}
	return
}

// Helper function to convert time.Time into microseconds since midnight
func toMicrosSinceMidnight(t time.Time) int64 {
	hour, min, sec := t.Clock()
	return int64(hour)*int64(time.Hour/time.Microsecond) +
		int64(min)*int64(time.Minute/time.Microsecond) +
		int64(sec)*int64(time.Second/time.Microsecond) +
		int64(t.Nanosecond())/int64(time.Microsecond)
}
