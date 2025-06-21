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
	"github.com/angryscorp/gophermart/internal/repository/orders/mapper"
	"github.com/angryscorp/gophermart/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orders struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

var _ repository.Orders = (*Orders)(nil)

func New(dsn string) (*Orders, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Orders{
		queries: db.New(pool),
		pool:    pool,
	}, nil
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
		UserID:     order.UserID,
		Status:     model.NewOrderStatus(order.Status),
		Accrual:    mapper.Mapper.NumericToFloat(order.Accrual),
		UploadedAt: order.UploadedAt.Time,
	}, nil
}

func (o Orders) CreateOrder(ctx context.Context, order model.Order) error {
	err := o.queries.CreateOrder(ctx, db.CreateOrderParams{
		Number: order.Number,
		UserID: order.UserID,
		Status: string(order.Status),
	})
	if err != nil {
		return err
	}

	return nil
}

func (o Orders) UpdateOrder(ctx context.Context, order model.Order) error {
	tx, err := o.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := o.queries.WithTx(tx)

	// Update order info
	if err := qtx.UpdateOrder(ctx, db.UpdateOrderParams{
		Status:  string(order.Status),
		Accrual: mapper.Mapper.FloatToNumeric(order.Accrual),
		Number:  order.Number,
	}); err != nil {
		return err
	}

	// Update balance
	if err := qtx.IncreaseBalance(ctx, order.Number); err != nil {
		return err
	}

	return nil
}

func (o Orders) AllOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	orders, err := o.queries.AllOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.Map(orders, func(item db.Order) model.Order {
		return model.Order{
			Number:     item.Number,
			UserID:     item.UserID,
			Status:     model.NewOrderStatus(item.Status),
			Accrual:    mapper.Mapper.NumericToFloat(item.Accrual),
			UploadedAt: item.UploadedAt.Time,
		}
	}), nil
}
