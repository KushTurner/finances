package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	dsn := os.Getenv("FINANCES_DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("FINANCES_DATABASE_URL environment variable is not set")
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}
