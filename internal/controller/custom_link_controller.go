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

var ErrCustomLinkController = "[CustomLinkController] Failed To Execute"

type CustomLinkControllerImpl struct {
	Service service.CustomLinkService
	Logger  *logger.Logger
}

func NewCustomLinkController(service service.CustomLinkService, logger *logger.Logger) CustomLinkController {
	return &CustomLinkControllerImpl{
		Service: service,
		Logger:  logger,
	}
}

func (controller *CustomLinkControllerImpl) CreateLink(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	var request web.CustomLinkCreateRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	customLinkResponse, errService := controller.Service.CreateLink(ctx, request, domainName, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success create new link",
			Data:    customLinkResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *CustomLinkControllerImpl) UpdateLink(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	var request web.CustomLinkUpdateRequest

	err := c.ShouldBindUri(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	customLinkResponse, errService := controller.Service.UpdateLink(ctx, request, domainName, jwtToken)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success update link",
			Data:    customLinkResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}

}

func (controller *CustomLinkControllerImpl) GetLink(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	var request web.CustomLinkGetRequest

	err := c.ShouldBindUri(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	customLinkResponse, errService := controller.Service.GetLink(ctx, request, domainName, jwtToken)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get link",
			Data:    customLinkResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}

func (controller *CustomLinkControllerImpl) GetAllLink(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	customLinksResponse, errService := controller.Service.GetAllLink(ctx, domainName, jwtToken)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get all custom link",
			Data:    customLinksResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}

}

func (controller *CustomLinkControllerImpl) GetAllThumbnail(c *gin.Context) {
	ctx := context.Background()
	thumbnailsResponse, errService := controller.Service.GetAllThumbnail(ctx)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get all default thumbnail",
			Data:    thumbnailsResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}

func (controller *CustomLinkControllerImpl) GetUserThumbnail(c *gin.Context) {
	ctx := context.Background()
	domainName := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	thumbnailsResponse, errService := controller.Service.GetUserThumbnail(ctx, domainName, jwtToken)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get all user custom thumbnail",
			Data:    thumbnailsResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}

}

func (controller *CustomLinkControllerImpl) UploadCustomThumbnail(c *gin.Context) {
	// TODO : Implement formfile validation
	// TODO : Fix bug, cant upload webp or png

	ctx := context.Background()
	file, _, err := c.Request.FormFile("image_data")
	domainName := c.Request.Host

	controller.Logger.PanicIfErr(err, ErrCustomLinkController)
	defer file.Close()

	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	controller.Logger.PanicIfErr(err, ErrCustomLinkController)

	thumbnailResponse, errService := controller.Service.UploadCustomThumbnail(ctx, buf.Bytes(), domainName, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success upload custom thumbnail",
			Data:    thumbnailResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *CustomLinkControllerImpl) CheckShortLinkAvaibility(c *gin.Context) {
	ctx := context.Background()
	var request web.CustomLinkCheckShortCodeAvaibilityRequest

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	errService := controller.Service.CheckShortLinkAvaibility(ctx, request)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "short link code can be used",
		}
		c.JSON(http.StatusOK, webResponse)
	}

}
func (controller *CustomLinkControllerImpl) RedirectLink(c *gin.Context) {
	ctx := context.Background()
	var request web.CustomLinkRedirectRequest

	err := c.ShouldBindUri(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	longLink, _, errService := controller.Service.RedirectLink(ctx, request)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		c.Redirect(http.StatusFound, longLink)
	}

}
