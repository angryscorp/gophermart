package orders

import (
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Orders struct {
	usecase usecase.Orders
}

var _ router.OrdersHandler = (*Orders)(nil)

func New(usecase usecase.Orders) Orders {
	return Orders{usecase: usecase}
}

func (r Orders) UploadOrder(c *gin.Context) {
	err := r.usecase.UploadOrder(c, "orderNumber", "username")
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, "UploadOrder")
}

func (r Orders) AllOrders(c *gin.Context) {
	orders, err := r.usecase.AllOrders(c)
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, orders)
}
