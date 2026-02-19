package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/config"
	"github.com/cenkalti/backoff"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	poolConfig.MaxConns = cfg.DB.MaxConns
	poolConfig.MinConns = cfg.DB.MinConns
	poolConfig.MaxConnLifetime = cfg.DB.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DB.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	if err := pingWithRetry(ctx, pool, cfg.DB.PingTimeout); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool, pingTimeout time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 200 * time.Millisecond
	bo.MaxInterval = 2 * time.Second
	bo.MaxElapsedTime = 12 * time.Second

	operation := func() error {
		pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			return fmt.Errorf("ping db: %w", err)
		}
		return nil
	}

	if err := backoff.Retry(operation, backoff.WithContext(bo, ctx)); err != nil {
		return fmt.Errorf("database is not reachable after retries: %w", err)
	}

	return nil
}
