package auth

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	usecase usecase.Auth
}

var _ router.AuthHandler = (*Auth)(nil)

type Request struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func New(usecase usecase.Auth) Auth {
	return Auth{usecase: usecase}
}

func (a Auth) SignIn(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, "Invalid request format")
		return
	}

	token, err := a.usecase.SignIn(c, req.Login, req.Password)
	if err != nil {
		c.JSON(401, "Credentials are invalid")
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(200, "SignIn Token: "+token+"")
}

func (a Auth) SignUp(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, "Invalid request format")
		return
	}

	token, err := a.usecase.SignUp(c, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserIsAlreadyExist) {
			c.JSON(409, "User is already exist")
		} else {
			c.JSON(500, "Something went wrong")
		}
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(200, "SignUp Token: "+token+"")
}
