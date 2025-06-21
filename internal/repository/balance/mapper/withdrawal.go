package mapper

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
)

type withdrawal struct{}

var Withdrawal = withdrawal{}

func (m withdrawal) ToDomainModel(row db.WithdrawalsRow) model.Withdrawal {
	v, _ := row.Withdrawn.Float64Value()
	return model.Withdrawal{
		Order:       row.OrderNumber,
		Sum:         v.Float64,
		ProcessedAt: row.ProcessedAt.Time,
	}
}
