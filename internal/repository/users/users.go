package users

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/repository/users/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Users struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

var _ repository.Users = (*Users)(nil)

func New(dsn string) (*Users, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Users{
		queries: db.New(pool),
		pool:    pool,
	}, nil
}

func (r Users) CreateUser(ctx context.Context, username, passwordHash string) (*uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := r.queries.WithTx(tx)

	isExist, err := qtx.CheckUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if isExist {
		return nil, model.ErrUserIsAlreadyExist
	}

	newID := uuid.New()
	if err := qtx.CreateUser(ctx, db.CreateUserParams{
		ID:           newID,
		Username:     username,
		PasswordHash: passwordHash,
	}); err != nil {
		return nil, errors.Wrap(model.ErrUnknownInternalError, err.Error())
	}

	err = qtx.CreateBalance(ctx, newID)
	if err != nil {
		return nil, errors.Wrap(model.ErrUnknownInternalError, err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &newID, nil
}

func (r Users) UserData(ctx context.Context, username string) (*model.User, error) {
	userData, err := r.queries.UserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:           userData.ID,
		Username:     userData.Username,
		PasswordHash: userData.PasswordHash,
	}, nil
}
