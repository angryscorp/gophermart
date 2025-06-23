package orders

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/utils/luhn"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type Orders struct {
	repository   repository.Orders
	requestChan  chan string
	responseChan chan *model.Accrual
	ctx          context.Context
	logger       zerolog.Logger
}

var _ usecase.Orders = (*Orders)(nil)

func New(
	repository repository.Orders,
	requestChan chan string,
	responseChan chan *model.Accrual,
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

func (o Orders) UploadOrder(ctx context.Context, orderNumber string, userID uuid.UUID) error {
	if orderNumber == "" {
		return usecase.ErrOrderNumberIsInvalid
	}

	numberIsValid := luhn.Validate(orderNumber)
	if !numberIsValid {
		return usecase.ErrOrderNumberIsInvalid
	}

	order, err := o.repository.OrderInfoForUpdate(ctx, orderNumber)
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to get order info")
		return model.ErrUnknownInternalError
	}

	if order != nil {
		if order.UserID == userID {
			return usecase.ErrOrderIsAlreadyUploaded
		} else {
			return usecase.ErrOrderWasUploadedAnotherUser
		}
	}

	err = o.repository.CreateOrder(ctx, model.NewOrder(orderNumber, userID))
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to create order")
		return model.ErrUnknownInternalError
	}

	o.requestChan <- orderNumber
	return nil
}

func (o Orders) AllOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	return o.repository.AllOrders(ctx, userID)
}

func (o Orders) listenResponses() {
	for resp := range o.responseChan {
		o.logger.Debug().Msgf("received response from accrual service: %+v", resp)
		order := newOrder(resp)
		if order != nil {
			o.logger.Debug().Msgf("updating order: %+v", order)
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

func newOrder(accrual *model.Accrual) *model.Order {
	switch accrual.Status {
	case model.AccrualStatusProcessed:
		return &model.Order{
			Number:  accrual.Order,
			Status:  model.OrderStatusProcessed,
			Accrual: accrual.Accrual,
		}

	case model.AccrualStatusInvalid:
		return &model.Order{
			Number: accrual.Order,
			Status: model.OrderStatusInvalid,
		}

	default:
		// the order is still processing
		return nil
	}
}
