package orders

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/common"
	"github.com/angryscorp/gophermart/internal/repository/orders/db"
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

func (o Orders) OrderInfo(ctx context.Context) model.Order {
	//TODO implement me
	panic("implement me")
}

func (o Orders) CreateOrder(ctx context.Context, order model.Order) error {
	//TODO implement me
	panic("implement me")
}

func (o Orders) UpdateOrder(ctx context.Context, order model.Order) error {
	//TODO implement me
	panic("implement me")
}
