package repository

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/google/uuid"
)

type Balance interface {
	Balance(ctx context.Context, userID uuid.UUID) (model.Balance, error)
	Withdraw(ctx context.Context, userID uuid.UUID, orderID string, amount int) error
	WithdrawalHistory(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
}
