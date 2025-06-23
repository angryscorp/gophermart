package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Number     string      `json:"number"`
	UserID     uuid.UUID   `json:"-"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func NewOrder(number string, userID uuid.UUID) Order {
	return Order{
		Number: number,
		UserID: userID,
		Status: OrderStatusNew,
	}
}
