package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/service"
)

type UserControllerImpl struct {
	Service service.UserService
	Logger  *logger.Logger
}

func NewUserController(service service.UserService, logger *logger.Logger) UserController {
	return &UserControllerImpl{
		Service: service,
		Logger:  logger,
	}
}

func (controller *UserControllerImpl) Register(c *gin.Context) {
	ctx := context.Background()

	var request web.UserRegisterRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	userResponse, errService := controller.Service.Register(ctx, request)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success create a new account",
			Data:    userResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *UserControllerImpl) Login(c *gin.Context) {
	ctx := context.Background()

	var request web.UserLoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	tokenResponse, errService := controller.Service.Login(ctx, request)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "login success",
			Data:    tokenResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}

func (controller *UserControllerImpl) ChangePassword(c *gin.Context) {
	ctx := context.Background()

	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	var request web.UserChangePasswordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	errService := controller.Service.ChangePassword(ctx, request, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success changed password",
		}
		c.JSON(http.StatusOK, webResponse)
	}
}
