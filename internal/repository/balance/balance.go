package balance

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/utils"
	"github.com/google/uuid"
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

	var current float64
	var withdrawn int

	if row.Balance.Valid {
		balanceFloat, _ := row.Balance.Float64Value()
		current = balanceFloat.Float64
	}

	if row.Withdrawn.Valid {
		withdrawnFloat, _ := row.Withdrawn.Float64Value()
		withdrawn = int(withdrawnFloat.Float64)
	}

	return model.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}

func (b Balance) Withdraw(ctx context.Context, userID uuid.UUID, orderID string, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (b Balance) WithdrawalHistory(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	rows, err := b.queries.Withdrawals(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get withdrawals")
	}

	return utils.Map(rows, func(row db.WithdrawalsRow) model.Withdrawal {
		v, _ := row.Withdrawn.Int64Value()
		return model.Withdrawal{
			Order:       row.OrderNumber,
			Sum:         int(v.Int64),
			ProcessedAt: row.ProcessedAt.Time,
		}
	}), nil
}
