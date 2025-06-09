package model

type Accrual struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual *int          `json:"accrual,omitempty"`
}
