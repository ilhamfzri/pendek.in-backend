package controller

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangePassword(c *gin.Context)
	Update(c *gin.Context)
	EmailVerification(c *gin.Context)
	GenerateToken(c *gin.Context)
	ChangeProfilePicture(c *gin.Context)
	Profile(c *gin.Context)
}

type SocialMediaLinkController interface {
	GetAllTypes(c *gin.Context)
	CreateLink(c *gin.Context)
	UpdateLink(c *gin.Context)
	GetAllLink(c *gin.Context)
	RedirectLink(c *gin.Context)
	GetLinkAnalytic(c *gin.Context)
	GetSummaryLinkAnalytic(c *gin.Context)
}

type CustomLinkController interface {
	CreateLink(c *gin.Context)
	UpdateLink(c *gin.Context)
	GetLink(c *gin.Context)
	GetAllLink(c *gin.Context)
	GetAllThumbnail(c *gin.Context)
	GetUserThumbnail(c *gin.Context)
	UploadCustomThumbnail(c *gin.Context)
	CheckShortLinkAvaibility(c *gin.Context)
	RedirectLink(c *gin.Context)
	GetLinkAnalytic(c *gin.Context)
	GetSummaryLinkAnalytic(c *gin.Context)
}
