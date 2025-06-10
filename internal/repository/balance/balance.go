package balance

import (
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/repository/balance/db"
	"github.com/angryscorp/gophermart/internal/repository/common"
)

type Balance struct {
	queries *db.Queries
}

var _ repository.Balance = (*Balance)(nil)

func New(dsn string) (*Balance, error) {
	pool, err := common.CreatePGXPool(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Balance{queries: db.New(pool)}, nil
}
