package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		if !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// TODO
		c.Set("username", "username")

		c.Next()
	}
}

func isValidToken(token string) bool {
	return token == "secret_token" // TODO
}
