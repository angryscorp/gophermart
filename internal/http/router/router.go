package router

import (
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/logger"
	"github.com/angryscorp/gophermart/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Router struct {
	engine         *gin.Engine
	authMiddleware gin.HandlerFunc
}

func New(
	zeroLogger zerolog.Logger,
	tokenValidator usecase.TokenValidator,
) Router {
	engine := gin.New()
	engine.
		Use(logger.New(zeroLogger)).
		Use(gin.Recovery())

	return Router{
		engine:         engine,
		authMiddleware: middleware.AuthValidation(tokenValidator, zeroLogger),
	}
}

func (r Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r Router) RegisterAuth(auth AuthHandler) {
	r.engine.POST("/api/user/register", auth.SignUp)
	r.engine.POST("/api/user/login", auth.SignIn)
}

func (r Router) RegisterOrders(orders OrdersHandler) {
	r.engine.POST("/api/user/orders", r.authMiddleware, orders.UploadOrder)
	r.engine.GET("/api/user/orders", r.authMiddleware, orders.AllOrders)
}

func (r Router) RegisterBalance(balance BalanceHandler) {
	r.engine.GET("/api/user/balance", r.authMiddleware, balance.Balance)
	r.engine.POST("/api/user/balance/withdraw", r.authMiddleware, balance.Withdraw)
	r.engine.GET("/api/user/withdrawals", r.authMiddleware, balance.AllWithdrawals)
}
