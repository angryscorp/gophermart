package repository

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/google/uuid"
)

type Orders interface {
	OrderInfoForUpdate(ctx context.Context, number string) (*model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	AllOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
}
