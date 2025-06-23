package usecase

import "github.com/angryscorp/gophermart/internal/domain/model"

type Accrual interface {
	Status(orderNumber string) (*model.Accrual, error)
}
