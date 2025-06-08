package handler

import (
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	usecase usecase.Auth
}

var _ router.AuthHandler = (*Auth)(nil)

func NewAuth(usecase usecase.Auth) Auth {
	return Auth{usecase: usecase}
}

func (a Auth) SignIn(c *gin.Context) {
	err := a.usecase.SignIn(c, "username", "password")
	if err != nil {
		c.JSON(401, "Unauthorized")
		return
	}
	c.JSON(200, "SignIn")
}

func (a Auth) SignUp(c *gin.Context) {
	err := a.usecase.SignUp(c, "username", "password")
	if err != nil {
		c.JSON(500, "Something went wrong")
		return
	}
	c.JSON(200, "SignUp")
}
