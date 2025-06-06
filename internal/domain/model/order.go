package model

import "time"

type Order struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
