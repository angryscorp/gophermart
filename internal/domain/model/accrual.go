package model

type Accrual struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual *float64      `json:"accrual,omitempty"`
}
