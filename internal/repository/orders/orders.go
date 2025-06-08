package orders

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/orders/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Orders struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var _ repository.Orders = (*Orders)(nil)

const (
	ctxTimeout = 10 * time.Second
)

func New(dsn string) (*Orders, error) {
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

	return &Orders{
		pool:    pool,
		queries: db.New(pool),
	}, nil
}
