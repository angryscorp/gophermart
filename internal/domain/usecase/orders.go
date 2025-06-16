package usecase

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

const (
	ErrOrderIsAlreadyUploaded      model.Error = "order has been already uploaded earlier"
	ErrOrderWasUploadedAnotherUser model.Error = "order was uploaded by another user"
	ErrOrderNumberIsInvalid        model.Error = "order number is invalid"
)

type Orders interface {
	UploadOrder(ctx context.Context, orderNumber, username string) error
	AllOrders(ctx context.Context) ([]model.Order, error)
}
