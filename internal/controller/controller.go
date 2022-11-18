package controller

import "github.com/gin-gonic/gin"

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangePassword(c *gin.Context)
	Update(c *gin.Context)
	EmailVerification(c *gin.Context)
	GenerateToken(c *gin.Context)
	ChangeProfilePicture(c *gin.Context)
}
