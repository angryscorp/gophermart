package repository

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

type Orders interface {
	OrderInfoForUpdate(ctx context.Context, number string) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	AllOrders(ctx context.Context) ([]model.Order, error)
}
