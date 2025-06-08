package balance

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Balance struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var _ repository.Balance = (*Balance)(nil)

const (
	ctxTimeout = 10 * time.Second
)

func New(dsn string) (*Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Balance{
		pool:    pool,
		queries: db.New(pool),
	}, nil
}
