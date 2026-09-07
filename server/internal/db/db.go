// Package db owns the Postgres connection and the one-time schema setup.
// There's no migration tool here on purpose: the schema is small enough that
// a single idempotent script, run at startup, is more legible than a
// migration framework for a project this size.
package db

import (
	"context"
	_ "embed"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed schema.sql
var schema string

// Connect opens a Postgres connection pool, verifies it's reachable, and
// applies schema.sql. The returned *sql.DB is ready to use.
func Connect(databaseURL string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if _, err := conn.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return conn, nil
}
