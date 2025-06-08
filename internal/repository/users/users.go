package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Users struct {
	pool *pgxpool.Pool
}

var _ repository.Users = (*Users)(nil)

const (
	ctxTimeout = 10 * time.Second
)

func New(dsn string) (*Users, error) {
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

	return &Users{
		pool: pool,
	}, nil
}

func (r Users) CreateUser(ctx context.Context, username, passwordHash string) error {
	fmt.Println("-->>CreateUser", username, passwordHash)

	const query = `
    INSERT INTO users (username, password_hash)
	VALUES ($1, $2)
	`

	_, err := r.pool.Exec(ctx, query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update metric: %w", err)
	}

	return nil
}

func (r Users) CheckUser(ctx context.Context, username, passwordHash string) error {
	fmt.Println("-->>CheckUser", username, passwordHash)

	const query = `
	SELECT 1 
	FROM users 
	WHERE 
		username = $1 
	  AND 
		password_hash = $2
	`

	var exists int
	err := r.pool.QueryRow(ctx, query, username, passwordHash).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("match not found")
	} else if err != nil {
		return fmt.Errorf("unknown error")
	}

	return nil
}
