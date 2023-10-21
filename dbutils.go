package main

import (
	"context"
	"os"

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
