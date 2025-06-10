package model

import "time"

type Order struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func NewOrder(number string) Order {
	return Order{
		Number:     number,
		Status:     OrderStatusNew,
		Accrual:    0,
		UploadedAt: time.Now(),
	}
}
