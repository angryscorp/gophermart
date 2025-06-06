package router

import (
	"github.com/angryscorp/gophermart/internal/http/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Router struct {
	engine *gin.Engine
}

func New(zeroLogger zerolog.Logger) Router {
	engine := gin.New()
	engine.
		Use(logger.New(zeroLogger)).
		Use(gin.Recovery())

	return Router{
		engine: engine,
	}
}

func (r Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r Router) RegisterAuth(auth AuthHandler) {
	r.engine.POST("/api/user/register", auth.SignUp)
	r.engine.POST("/api/user/login", auth.SignIn)
}
