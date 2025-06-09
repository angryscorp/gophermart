package orders

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"time"
)

type Orders struct {
	repository   repository.Orders
	requestChan  chan string
	responseChan chan model.Accrual
}

var _ usecase.Orders = (*Orders)(nil)

func New(
	repository repository.Orders,
	requestChan chan string,
	responseChan chan model.Accrual,
) Orders {
	orders := Orders{
		repository:   repository,
		requestChan:  requestChan,
		responseChan: responseChan,
	}

	go orders.listenResponses()

	return orders
}

func (o Orders) UploadOrder(ctx context.Context, orderNumber string) error {
	o.requestChan <- orderNumber
	return nil // TODO
}

func (o Orders) AllOrders(ctx context.Context) ([]model.Order, error) {
	return []model.Order{{Number: "234", Status: model.OrderStatusProcessing, Accrual: 123, UploadedAt: time.Now()}}, nil // TODO
}

func (o Orders) listenResponses() {
	for resp := range o.responseChan {
		fmt.Println(resp)
	}
}
