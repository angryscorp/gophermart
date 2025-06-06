package router

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
}
