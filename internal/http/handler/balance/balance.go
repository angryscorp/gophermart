package balance

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Balance struct {
	usecase usecase.Balance
}

var _ router.BalanceHandler = (*Balance)(nil)

func New(usecase usecase.Balance) Balance {
	return Balance{
		usecase: usecase,
	}
}

func (b Balance) Balance(c *gin.Context) {
	balance, err := b.usecase.Balance(c)
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, balance)
}

func (b Balance) Withdraw(c *gin.Context) {
	err := b.usecase.Withdraw(c, model.WithdrawalRequest{
		Order: "123",
		Sum:   123,
	})
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, "Withdrawal")
}

func (b Balance) AllWithdrawals(c *gin.Context) {
	withdrawals, err := b.usecase.AllWithdrawals(c)
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, withdrawals)
}
