package mapper

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type balance struct{}

var Balance = balance{}

func (m balance) ToDomainModel(row db.BalanceRow) model.Balance {
	return model.Balance{
		Current:   m.NumericToFloat(row.Balance),
		Withdrawn: m.NumericToInt(row.Withdrawn),
	}
}

func (m balance) NumericToFloat(value pgtype.Numeric) float64 {
	v, _ := value.Float64Value()
	return v.Float64
}

func (m balance) NumericToInt(value pgtype.Numeric) int {
	v, _ := value.Int64Value()
	return int(v.Int64)
}

func (m balance) IntToNumeric(value int) pgtype.Numeric {
	v := pgtype.Numeric{}
	_ = v.Scan(value)
	return v
}
