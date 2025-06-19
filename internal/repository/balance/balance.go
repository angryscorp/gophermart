package balance

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/google/uuid"
)

type Balance struct {
	queries *db.Queries
}

var _ repository.Balance = (*Balance)(nil)

func New(dsn string) (*Balance, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Balance{queries: db.New(pool)}, nil
}

func (b Balance) Balance(ctx context.Context, userID uuid.UUID) (repository.Balance, error) {
	//TODO implement me
	panic("implement me")
}

func (b Balance) Withdraw(ctx context.Context, userID uuid.UUID, orderID string, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (b Balance) WithdrawalHistory(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	//TODO implement me
	panic("implement me")
}
