package repository

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

type Orders interface {
	OrderInfo(ctx context.Context) model.Order
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	AllOrders(ctx context.Context) ([]model.Order, error)
}
