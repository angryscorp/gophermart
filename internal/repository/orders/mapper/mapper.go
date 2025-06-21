package mapper

import "github.com/jackc/pgx/v5/pgtype"

type mapper struct{}

var Mapper = mapper{}

func (m mapper) NumericToFloat(value pgtype.Numeric) *float64 {
	if !value.Valid {
		return nil
	}
	v, err := value.Float64Value()
	if err != nil {
		return nil
	}
	return &v.Float64
}

func (m mapper) FloatToNumeric(value *float64) pgtype.Numeric {
	v := pgtype.Numeric{}
	_ = v.Scan(value)
	return v
}
