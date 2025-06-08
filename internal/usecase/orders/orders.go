package orders

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"time"
)

type Orders struct {
	repository repository.Orders
}

var _ usecase.Orders = (*Orders)(nil)

func New(repository repository.Orders) Orders {
	return Orders{
		repository: repository,
	}
}

func (o Orders) UploadOrder(ctx context.Context, orderNumber string) error {
	return nil // TODO
}

func (o Orders) AllOrders(ctx context.Context) ([]model.Order, error) {
	return []model.Order{{Number: "234", Status: model.OrderStatusProcessing, Accrual: 123, UploadedAt: time.Now()}}, nil // TODO
}
