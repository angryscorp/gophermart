package balance

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/handler/common"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Balance struct {
	usecase usecase.Balance
	logger  zerolog.Logger
}

var _ router.BalanceHandler = (*Balance)(nil)

type Request struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

func New(
	usecase usecase.Balance,
	logger zerolog.Logger,
) Balance {
	return Balance{
		usecase: usecase,
		logger:  logger,
	}
}

func (b Balance) Balance(c *gin.Context) {
	b.logger.Debug().Msg("Handler Balance")

	userID, err := common.GetUserID(c)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to get user ID")
		c.JSON(500, "Internal server error")
		return
	}

	b.logger.Debug().Str("user_id", userID.String()).Msg("User ID")

	balance, err := b.usecase.Balance(c, *userID)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to get balance")
		c.JSON(500, "Something went wrong")
		return
	}

	c.JSON(200, balance)
}

func (b Balance) Withdraw(c *gin.Context) {
	b.logger.Debug().Msg("Handler Withdraw")

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Debug().Err(err).Msg("Failed to bind request")
		c.JSON(400, "Invalid request format")
		return
	}

	b.logger.Debug().Interface("request", req).Msg("Request")

	userID, err := common.GetUserID(c)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to get user ID")
		c.JSON(500, "Internal server error")
		return
	}

	b.logger.Debug().Str("user_id", userID.String()).Msg("User ID")

	err = b.usecase.Withdraw(c, *userID, req.Order, req.Sum)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to withdraw")
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
	b.logger.Debug().Msg("Handler AllWithdrawals")

	userID, err := common.GetUserID(c)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to get user ID")
		c.JSON(500, "Internal server error")
		return
	}

	b.logger.Debug().Str("user_id", userID.String()).Msg("User ID")

	withdrawals, err := b.usecase.AllWithdrawals(c, *userID)
	if err != nil {
		b.logger.Debug().Err(err).Msg("Failed to get withdrawals")
		c.JSON(500, "Something went wrong")
		return
	}

	b.logger.Debug().Int("withdrawals_count", len(withdrawals)).Msg("Withdrawals count")

	if len(withdrawals) == 0 {
		c.JSON(204, "No withdrawals")
		return
	}

	c.JSON(200, withdrawals)
}
