package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"k8s.io/klog/v2"
)

//go:embed schema.sql
var schemaDDL string

// SchemaSQL returns the embedded schema DDL (CREATE TABLE / indexes / constraints).
func SchemaSQL() string {
	return schemaDDL
}

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return pool, nil
}

func Bootstrap(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, schemaDDL); err != nil {
		return fmt.Errorf("schema bootstrap: %w", err)
	}
	klog.Info("schema bootstrap done")
	return nil
}

func CleanupData(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `TRUNCATE match_attendees, payments, matches, players RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("cleanup data: %w", err)
	}
	return nil
}
