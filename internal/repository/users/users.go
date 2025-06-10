package users

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/repository/users/db"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Users struct {
	queries *db.Queries
}

var _ repository.Users = (*Users)(nil)

func New(dsn string) (*Users, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Users{queries: db.New(pool)}, nil
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
