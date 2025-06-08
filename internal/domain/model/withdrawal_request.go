package model

type WithdrawalRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}
