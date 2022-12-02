package controller

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/service"
)

type UserControllerImpl struct {
	Service            service.UserService
	SocialMediaService service.SocialMediaLinkService
	CustomLinkService  service.CustomLinkService
	Logger             *logger.Logger
}

var ErrUserController = "[UserController] Failed To Execute"

func NewUserController(service service.UserService, socialMediaService service.SocialMediaLinkService, costumLinkService service.CustomLinkService, logger *logger.Logger) UserController {
	return &UserControllerImpl{
		Service:            service,
		SocialMediaService: socialMediaService,
		CustomLinkService:  costumLinkService,
		Logger:             logger,
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

	// TODO : integrate with email services to send verification code to user email
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

func (controller *UserControllerImpl) Update(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	var request web.UserUpdateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	userResponse, errService := controller.Service.Update(ctx, request, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success update user info",
			Data:    userResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *UserControllerImpl) EmailVerification(c *gin.Context) {
	ctx := context.Background()
	var request web.UserEmailVerificationRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	userResponse, errService := controller.Service.EmailVerification(ctx, request)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "successfully verified the email",
			Data:    userResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *UserControllerImpl) GenerateToken(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	tokenResponse, errService := controller.Service.GenerateToken(ctx, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success create new request token",
			Data:    tokenResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}

}

func (controller *UserControllerImpl) ChangeProfilePicture(c *gin.Context) {
	ctx := context.Background()
	file, _, err := c.Request.FormFile("image_data")
	controller.Logger.PanicIfErr(err, ErrUserController)
	defer file.Close()

	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	controller.Logger.PanicIfErr(err, ErrUserController)

	errService := controller.Service.ChangeProfilePicture(ctx, buf.Bytes(), jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success change profile picture",
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *UserControllerImpl) Profile(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	var request web.UserProfileRequest

	err := c.ShouldBindUri(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	userResponse := controller.Service.GetProfileData(ctx, request)
	if userResponse.ID == "" {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: "account not found",
		}
		c.JSON(http.StatusBadRequest, webResponse)
	}

	userProfileResponse := web.UserProfileResponse{}
	userProfileResponse.Username = userResponse.Username
	userProfileResponse.FullName = userResponse.FullName
	userProfileResponse.Bio = userResponse.Bio

	controller.SocialMediaService.GetAllLinkProfile(ctx, domainName, userResponse.ID, userResponse.Username)
	socialMediaResponse := controller.SocialMediaService.GetAllLinkProfile(ctx, domainName, userResponse.ID, userResponse.Username)
	customLinkResponse := controller.CustomLinkService.GetAllLinkProfile(ctx, domainName, userResponse.ID, userResponse.Username)

	if socialMediaResponse != nil {
		userProfileResponse.SocialMedia = socialMediaResponse
	}

	if customLinkResponse != nil {
		userProfileResponse.Link = customLinkResponse
	}

	if userResponse.ProfilePic != "" {
		userProfileResponse.ProfilePic = helper.GetProfilePictureUrl(domainName, userResponse.ProfilePic)
	}

	webResponse := web.WebResponseSuccess{
		Status:  "success",
		Message: "success get account profile",
		Data:    userProfileResponse,
	}
	c.JSON(http.StatusOK, webResponse)

}

func (controller *UserControllerImpl) GetCurrentProfile(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	domainName := c.Request.Host

	userResponse, errService := controller.Service.GetCurrentProfile(ctx, jwtToken)

	if userResponse.ProfilePic != "" {
		userResponse.ProfilePic = helper.GetProfilePictureUrl(domainName, userResponse.ProfilePic)
	}

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get current profile",
			Data:    userResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}
