package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Number     string      `json:"number"`
	UserId     uuid.UUID   `json:"-"`
	Status     OrderStatus `json:"status"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func NewOrder(number string, userId uuid.UUID) Order {
	return Order{
		Number: number,
		UserId: userId,
		Status: OrderStatusNew,
	}
}
