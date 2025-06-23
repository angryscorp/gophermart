package auth

import (
	"errors"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/angryscorp/gophermart/internal/http/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Auth struct {
	usecase usecase.Auth
	logger  zerolog.Logger
}

var _ router.AuthHandler = (*Auth)(nil)

type Request struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func New(
	usecase usecase.Auth,
	logger zerolog.Logger,
) Auth {
	return Auth{
		usecase: usecase,
		logger:  logger,
	}
}

func (a Auth) SignIn(c *gin.Context) {
	a.logger.Debug().Msg("Handler SignIn")

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Debug().Err(err).Msg("Failed to bind request")
		c.JSON(400, "Invalid request format")
		return
	}

	a.logger.Debug().Interface("request", req).Msg("Request")

	token, err := a.usecase.SignIn(c, req.Login, req.Password)
	if err != nil {
		a.logger.Debug().Err(err).Msg("Failed to sign in")
		c.JSON(401, "Credentials are invalid")
		return
	}

	a.logger.Debug().Str("token", token).Msg("Token")

	c.Header("Authorization", "Bearer "+token)
	c.JSON(200, "SignIn Token: "+token+"")
}

func (a Auth) SignUp(c *gin.Context) {
	a.logger.Debug().Msg("Handler SignUp")

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Debug().Err(err).Msg("Failed to bind request")
		c.JSON(400, "Invalid request format")
		return
	}

	a.logger.Debug().Interface("request", req).Msg("Request")

	token, err := a.usecase.SignUp(c, req.Login, req.Password)
	if err != nil {
		a.logger.Debug().Err(err).Msg("Failed to sign up")
		if errors.Is(err, model.ErrUserIsAlreadyExist) {
			c.JSON(409, "User is already exist")
		} else {
			c.JSON(500, "Something went wrong")
		}
		return
	}

	a.logger.Debug().Str("token", token).Msg("Token")

	c.Header("Authorization", "Bearer "+token)
	c.JSON(200, "SignUp Token: "+token+"")
}
