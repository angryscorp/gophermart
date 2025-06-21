package orders

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/handler/common"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"strings"
)

type Orders struct {
	usecase usecase.Orders
	logger  zerolog.Logger
}

var _ router.OrdersHandler = (*Orders)(nil)

func New(usecase usecase.Orders, logger zerolog.Logger) Orders {
	return Orders{
		usecase: usecase,
		logger:  logger,
	}
}

func (r Orders) UploadOrder(c *gin.Context) {
	r.logger.Debug().Msg("Handler UploadOrder")

	// Order number
	orderNumberBytes, err := c.GetRawData()
	if err != nil {
		r.logger.Debug().Err(err).Msg("Failed to get order number")
		c.JSON(400, "Invalid request format")
		return
	}

	orderNumber := strings.TrimSpace(string(orderNumberBytes))
	r.logger.Debug().Str("order_number", orderNumber).Msg("Order number")

	// userID
	userID, err := common.GetUserID(c)
	if err != nil {
		r.logger.Debug().Err(err).Msg("Failed to get user ID")
		c.JSON(500, "Internal server error")
		return
	}

	r.logger.Debug().Str("user_id", userID.String()).Msg("User ID")

	// Mail logic
	err = r.usecase.UploadOrder(c, orderNumber, *userID)
	if err != nil {
		r.logger.Debug().Err(err).Msg("Failed to upload order")
		switch {
		case errors.Is(err, usecase.ErrOrderIsAlreadyUploaded):
			c.JSON(200, "Order is already uploaded")

		case errors.Is(err, usecase.ErrOrderNumberIsInvalid):
			c.JSON(422, "Order number is invalid")

		case errors.Is(err, usecase.ErrOrderWasUploadedAnotherUser):
			c.JSON(409, "Order was uploaded by another user")

		default:
			c.JSON(500, "Something went wrong")
		}

		return
	}

	c.JSON(202, "Order is uploaded")
}

func (r Orders) AllOrders(c *gin.Context) {
	r.logger.Debug().Msg("Handler AllOrders")

	// userID
	userID, err := common.GetUserID(c)
	if err != nil {
		r.logger.Debug().Err(err).Msg("Failed to get user ID")
		c.JSON(500, "Internal server error")
		return
	}
	r.logger.Debug().Str("user_id", userID.String()).Msg("User ID")

	// Main logic
	orders, err := r.usecase.AllOrders(c, *userID)
	if err != nil {
		r.logger.Debug().Err(err).Msg("Failed to get orders")
		c.JSON(500, "Internal server error")
		return
	}
	r.logger.Debug().Int("orders_count", len(orders)).Msg("Orders count")

	if len(orders) == 0 {
		c.JSON(204, "No orders")
		return
	}

	c.JSON(200, orders)
}
