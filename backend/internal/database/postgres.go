package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

// PostgresRepository implements the ReceiptRepository interface using PostgreSQL
type PostgresRepository struct {
	Pool *pgxpool.Pool
}

// Ensure PostgresRepository implements ReceiptRepository interface
var _ ReceiptRepository = (*PostgresRepository)(nil)

func NewPostgres(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
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

	repo := &PostgresRepository{Pool: pool}

	// Initialize schema
	if err := repo.InitSchema(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *PostgresRepository) InitSchema(ctx context.Context) error {
	_, err := r.Pool.Exec(ctx, schemaSQL)
	return err
}

func (r *PostgresRepository) Close() {
	if r.Pool != nil {
		r.Pool.Close()
	}
}
