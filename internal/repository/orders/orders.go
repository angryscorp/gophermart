package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/repository/orders/db"
	"github.com/angryscorp/gophermart/internal/utils"
)

type Orders struct {
	queries *db.Queries
}

var _ repository.Orders = (*Orders)(nil)

func New(dsn string) (*Orders, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Orders{queries: db.New(pool)}, nil
}

func (o Orders) OrderInfoForUpdate(ctx context.Context, number string) (*model.Order, error) {
	order, err := o.queries.GetOrderForUpdate(ctx, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.Order{
		Number:     order.Number,
		Username:   order.Username,
		Status:     model.NewOrderStatus(order.Status),
		Accrual:    int(order.Accrual),
		UploadedAt: order.UploadedAt.Time,
	}, nil
}

func (o Orders) CreateOrder(ctx context.Context, order model.Order) error {
	err := o.queries.CreateOrder(ctx, db.CreateOrderParams{
		Number:   order.Number,
		Username: order.Username,
		Status:   string(order.Status),
	})
	if err != nil {
		return err
	}

	return nil
}

func (o Orders) UpdateOrder(ctx context.Context, order model.Order) error {
	if err := o.queries.UpdateOrder(ctx, db.UpdateOrderParams{
		Status:   string(order.Status),
		Accrual:  int32(order.Accrual),
		Number:   order.Number,
		Username: order.Username,
	}); err != nil {
		return err
	}

	return nil
}

func (o Orders) AllOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := o.queries.AllOrders(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Map(orders, func(item db.Order) model.Order {
		return model.Order{
			Number:     item.Number,
			Username:   item.Username,
			Status:     model.NewOrderStatus(item.Status),
			Accrual:    int(item.Accrual),
			UploadedAt: item.UploadedAt.Time,
		}
	}), nil
}
