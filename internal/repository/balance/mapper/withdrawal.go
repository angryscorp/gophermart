package mapper

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
)

type withdrawal struct{}

var Withdrawal = withdrawal{}

func (m withdrawal) ToDomainModel(row db.WithdrawalsRow) model.Withdrawal {
	v, _ := row.Withdrawn.Int64Value()
	return model.Withdrawal{
		Order:       row.OrderNumber,
		Sum:         int(v.Int64),
		ProcessedAt: row.ProcessedAt.Time,
	}
}
