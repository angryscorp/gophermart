package model

import "time"

type Order struct {
	Number     string      `json:"number"`
	Username   string      `json:"username"`
	Status     OrderStatus `json:"status"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func NewOrder(number, username string) Order {
	return Order{
		Number:   number,
		Username: username,
		Status:   OrderStatusNew,
	}
}
