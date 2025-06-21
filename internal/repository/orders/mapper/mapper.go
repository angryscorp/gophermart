package mapper

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

type mapper struct{}

var Mapper = mapper{}

func (m mapper) NumericToFloat(value pgtype.Numeric) *float64 {
	v, _ := value.Float64Value()
	return &v.Float64
}

func (m mapper) FloatToNumeric(value *float64) pgtype.Numeric {
	v := pgtype.Numeric{}
	_ = v.Scan(strconv.FormatFloat(*value, 'f', 2, 64))
	return v

}
