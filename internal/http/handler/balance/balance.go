package balance

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/handler/common"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Balance struct {
	usecase usecase.Balance
}

var _ router.BalanceHandler = (*Balance)(nil)

type Request struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

func New(usecase usecase.Balance) Balance {
	return Balance{
		usecase: usecase,
	}
}

func (b Balance) Balance(c *gin.Context) {
	userID, err := common.GetUserID(c)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	balance, err := b.usecase.Balance(c, *userID)
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}

	c.JSON(200, balance)
}

func (b Balance) Withdraw(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, "Invalid request format")
		return
	}

	userID, err := common.GetUserID(c)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	err = b.usecase.Withdraw(c, *userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrOrderNumberIsInvalid):
			c.JSON(422, "Order number is invalid")

		case errors.Is(err, model.ErrUnsufficientFunds):
			c.JSON(402, "Unsufficient funds")

		default:
			c.JSON(500, "Something went wrong")
		}

		return
	}

	c.JSON(200, "Withdrawal")
}

func (b Balance) AllWithdrawals(c *gin.Context) {
	userID, err := common.GetUserID(c)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	withdrawals, err := b.usecase.AllWithdrawals(c, *userID)
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}

	if len(withdrawals) == 0 {
		c.JSON(204, "No withdrawals")
		return
	}
	
	c.JSON(200, withdrawals)
}
