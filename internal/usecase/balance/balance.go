package balance

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (b Balance) Balance(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	return b.repository.Balance(ctx, userID)
}

func (b Balance) Withdraw(ctx context.Context, userID uuid.UUID, orderNumber string, amount float64) error {
	if orderNumber == "" {
		return usecase.ErrOrderNumberIsInvalid
	}

	numberIsValid := utils.CheckLuhn(orderNumber)
	if !numberIsValid {
		return usecase.ErrOrderNumberIsInvalid
	}

	err := b.repository.Withdraw(ctx, userID, orderNumber, amount)
	if err != nil {
		return errors.Wrap(err, "failed to withdraw")
	}

	return nil
}

func (b Balance) AllWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	return b.repository.WithdrawalHistory(ctx, userID)
}
