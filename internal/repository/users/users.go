package users

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/users/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type Users struct {
	pool    *pgxpool.Pool
	queries *db.Queries
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
		pool:    pool,
		queries: db.New(pool),
	}, nil
}

func (r Users) CreateUser(ctx context.Context, username, passwordHash string) error {
	if err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	}); err != nil {
		return errors.Wrap(repository.ErrInvalidCredentials, err.Error())
	}

	return nil
}

func (r Users) CheckUser(ctx context.Context, username, passwordHash string) error {
	_, err := r.queries.CheckUser(ctx, db.CheckUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrInvalidCredentials
	} else if err != nil {
		return errors.Wrap(repository.ErrUnknownInternalError, err.Error())
	}

	return nil
}
