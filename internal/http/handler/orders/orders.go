package orders

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/handler/common"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
	"strings"
)

type Orders struct {
	usecase usecase.Orders
}

var _ router.OrdersHandler = (*Orders)(nil)

func New(usecase usecase.Orders) Orders {
	return Orders{usecase: usecase}
}

func (r Orders) UploadOrder(c *gin.Context) {
	// Order number
	orderNumberBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(400, "Invalid request format")
		return
	}

	orderNumber := strings.TrimSpace(string(orderNumberBytes))

	// userID
	userID, err := common.GetUserID(c)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	// Mail logic
	err = r.usecase.UploadOrder(c, orderNumber, *userID)

	switch {
	case errors.Is(err, usecase.ErrOrderIsAlreadyUploaded):
		c.JSON(200, "Order is already uploaded")

	case errors.Is(err, usecase.ErrOrderNumberIsInvalid):
		c.JSON(422, "Order number is invalid")

	case errors.Is(err, usecase.ErrOrderWasUploadedAnotherUser):
		c.JSON(409, "Order was uploaded by another user")

	case err == nil:
		c.JSON(202, "Order is uploaded")

	default:
		c.JSON(500, "Something went wrong")
	}
}

func (r Orders) AllOrders(c *gin.Context) {
	// userID
	userID, err := common.GetUserID(c)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	// Main logic
	orders, err := r.usecase.AllOrders(c, *userID)
	if err != nil {
		c.JSON(500, "Internal server error")
		return
	}

	if len(orders) == 0 {
		c.JSON(204, "No orders")
		return
	}

	c.JSON(200, orders)
}
