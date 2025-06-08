package usecase

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

type Orders interface {
	UploadOrder(ctx context.Context, orderNumber string) error
	AllOrders(ctx context.Context) ([]model.Order, error)
}
