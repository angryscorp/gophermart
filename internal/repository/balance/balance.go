package balance

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/angryscorp/gophermart/internal/repository/balance/mapper"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Balance struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

var _ repository.Balance = (*Balance)(nil)

func New(dsn string) (*Balance, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Balance{
		queries: db.New(pool),
		pool:    pool,
	}, nil
}

func (b Balance) Balance(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	row, err := b.queries.Balance(ctx, userID)
	if err != nil {
		return model.Balance{}, errors.Wrap(err, "failed to get balance")
	}

	return mapper.Balance.ToDomainModel(row), nil
}

func (b Balance) Withdraw(ctx context.Context, userID uuid.UUID, orderID string, amount float64) error {
	tx, err := b.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := b.queries.WithTx(tx)

	// Check available balance
	row, err := qtx.CheckBalanceForUpdate(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to check balance")
	}

	currentBalance := mapper.Balance.NumericToFloat(row.Balance)
	if currentBalance < amount {
		return model.ErrUnsufficientFunds
	}

	// New Balance
	newBalance := mapper.Balance.FloatToNumeric(currentBalance - amount)
	newWithdrawn := mapper.Balance.FloatToNumeric(mapper.Balance.NumericToFloat(row.Withdrawn) + amount)
	err = qtx.UpdateBalance(ctx, db.UpdateBalanceParams{
		Balance:   newBalance,
		Withdrawn: newWithdrawn,
		UserID:    userID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to update balance")
	}

	// Update history
	err = qtx.AddWithdrawal(ctx, db.AddWithdrawalParams{
		OrderNumber: orderID,
		UserID:      userID,
		Withdrawn:   mapper.Balance.FloatToNumeric(amount),
	})
	if err != nil {
		return errors.Wrap(err, "failed to add withdrawal")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (b Balance) WithdrawalHistory(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	rows, err := b.queries.Withdrawals(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get withdrawals")
	}

	return utils.Map(rows, func(row db.WithdrawalsRow) model.Withdrawal {
		return mapper.Withdrawal.ToDomainModel(row)
	}), nil
}
