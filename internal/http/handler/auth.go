package handler

import (
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Auth struct {
}

var _ router.AuthHandler = (*Auth)(nil)

func NewAuth() Auth {
	return Auth{}
}

func (a Auth) SignIn(c *gin.Context) {
	c.JSON(200, "SignIn")
}

func (a Auth) SignUp(c *gin.Context) {
	c.JSON(200, "SignUp")
}
