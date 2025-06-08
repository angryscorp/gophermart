package usecase

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

type Balance interface {
	Balance(ctx context.Context) (model.Balance, error)
	Withdraw(ctx context.Context, req model.WithdrawalRequest) error
	AllWithdrawals(ctx context.Context) ([]model.Withdrawal, error)
}
