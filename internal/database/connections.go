package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dbURL)
}
