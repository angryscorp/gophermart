package middleware

import (
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthValidation(
	tokenValidator usecase.TokenValidator,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := tokenValidator.Validate(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", token.UserID)
		c.Next()
	}
}
