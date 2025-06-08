package balance

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"time"
)

type Balance struct {
	repository repository.Balance
}

var _ usecase.Balance = (*Balance)(nil)

func New(repository repository.Balance) Balance {
	return Balance{
		repository: repository,
	}
}

func (b Balance) Balance(ctx context.Context) (model.Balance, error) {
	return model.Balance{Current: 123, Withdrawn: 45}, nil
}

func (b Balance) Withdraw(ctx context.Context, req model.WithdrawalRequest) error {
	return nil
}

func (b Balance) AllWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	return []model.Withdrawal{{Order: "123", Sum: 456, ProcessedAt: time.Now()}}, nil
}
