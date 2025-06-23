package usecase

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/google/uuid"
)

type Balance interface {
	Balance(ctx context.Context, userID uuid.UUID) (model.Balance, error)
	Withdraw(ctx context.Context, userID uuid.UUID, orderNumber string, amount float64) error
	AllWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
}
