package middleware

import (
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

func AuthValidation(
	tokenValidator usecase.TokenValidator,
	logger zerolog.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug().Msg("Middleware AuthValidation")

		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			logger.Debug().Msg("Authorization token is empty")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		logger.Debug().Msgf("Authorization token: %s", tokenStr)

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := tokenValidator.Validate(tokenStr)
		if err != nil {
			logger.Debug().Err(err).Msg("Failed to validate token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		logger.Debug().Interface("token", token).Msg("Token")

		c.Set("userID", token.UserID)
		c.Next()
	}
}
