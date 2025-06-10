package orders

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/rs/zerolog"
	"time"
)

type Orders struct {
	repository   repository.Orders
	requestChan  chan string
	responseChan chan model.Accrual
	ctx          context.Context
	logger       zerolog.Logger
}

var _ usecase.Orders = (*Orders)(nil)

func New(
	repository repository.Orders,
	requestChan chan string,
	responseChan chan model.Accrual,
	logger zerolog.Logger,
) Orders {
	orders := Orders{
		repository:   repository,
		requestChan:  requestChan,
		responseChan: responseChan,
		ctx:          context.Background(),
		logger:       logger,
	}

	go orders.listenResponses()

	return orders
}

func (o Orders) UploadOrder(ctx context.Context, orderNumber string) error {
	err := o.repository.CreateOrder(ctx, model.NewOrder(orderNumber))
	if err != nil {
		return err
	}

	o.requestChan <- orderNumber
	return nil
}

func (o Orders) AllOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := o.repository.AllOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o Orders) listenResponses() {
	for resp := range o.responseChan {
		order := newOrder(resp)
		if order != nil {
			err := o.repository.UpdateOrder(o.ctx, *order)
			if err != nil {
				o.logger.Error().Err(err).Msg("failed to update order")
			}
			go func() {
				time.Sleep(5 * time.Second)
				o.responseChan <- resp
			}()
		}
	}
}

func newOrder(accrual model.Accrual) *model.Order {
	switch accrual.Status {
	case model.AccrualStatusProcessed:
		accrualValue := 0
		if accrual.Accrual != nil {
			accrualValue = *accrual.Accrual
		}
		return &model.Order{
			Number:  accrual.Order,
			Status:  model.OrderStatusProcessed,
			Accrual: accrualValue,
		}

	case model.AccrualStatusInvalid:
		return &model.Order{
			Number:  accrual.Order,
			Status:  model.OrderStatusInvalid,
			Accrual: 0,
		}

	default:
		// the order is still processing
		return nil
	}
}
