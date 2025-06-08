package router

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
}

type OrdersHandler interface {
	UploadOrder(c *gin.Context)
	AllOrders(c *gin.Context)
}

type BalanceHandler interface {
	Balance(c *gin.Context)
	Withdraw(c *gin.Context)
	AllWithdrawals(c *gin.Context)
}
