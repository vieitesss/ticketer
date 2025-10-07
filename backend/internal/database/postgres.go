package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

type DB struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{Pool: pool}

	// Initialize schema
	if err := db.InitSchema(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func (db *DB) InitSchema(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, schemaSQL)
	return err
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
