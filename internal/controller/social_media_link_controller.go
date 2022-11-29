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

type SocialMediaLinkControllerImpl struct {
	Service         service.SocialMediaLinkService
	AnalyticService service.SocialMediaAnalytic
	Logger          *logger.Logger
}

func NewSocialMediaLink(service service.SocialMediaLinkService, analyticService service.SocialMediaAnalytic, logger *logger.Logger) SocialMediaLinkController {
	return &SocialMediaLinkControllerImpl{
		Service:         service,
		AnalyticService: analyticService,
		Logger:          logger,
	}
}

func (controller *SocialMediaLinkControllerImpl) GetAllTypes(c *gin.Context) {
	ctx := context.Background()
	socialMediaTypesResponse, errService := controller.Service.GetAllTypes(ctx)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get all social media types",
			Data:    socialMediaTypesResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}

func (controller *SocialMediaLinkControllerImpl) CreateLink(c *gin.Context) {
	ctx := context.Background()
	host := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	var request web.SocialMediaLinkCreateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	socialMediaLinkResponse, errService := controller.Service.CreateLink(ctx, request, host, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success create social media link",
			Data:    socialMediaLinkResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *SocialMediaLinkControllerImpl) UpdateLink(c *gin.Context) {
	ctx := context.Background()
	host := c.Request.Host
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	var request web.SocialMediaLinkUpdateRequest

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

	socialMediaLinkResponse, errService := controller.Service.UpdateLink(ctx, request, host, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success update social media link",
			Data:    socialMediaLinkResponse,
		}
		c.JSON(http.StatusCreated, webResponse)
	}
}

func (controller *SocialMediaLinkControllerImpl) GetAllLink(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	host := c.Request.Host

	socialMediaTypesResponse, errService := controller.Service.GetAllLink(ctx, host, jwtToken)
	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get all social media types",
			Data:    socialMediaTypesResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}
}

func (controller *SocialMediaLinkControllerImpl) RedirectLink(c *gin.Context) {
	ctx := context.Background()
	var request web.SocialMediaLinkRedirectRequest

	err := c.ShouldBindUri(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}
	redirectLink, socialMediaLinkID, errService := controller.Service.RedirectLink(ctx, request)

	if errService == nil {
		requstSaveInteraction := web.SocialMediaAnalyticInteractionRequest{
			ClientIP:          c.ClientIP(),
			UserAgent:         c.Request.Header.Get("User-Agent"),
			SocialMediaLinkID: socialMediaLinkID,
		}
		_ = controller.AnalyticService.SaveInteraction(ctx, requstSaveInteraction)
	}

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		c.Redirect(http.StatusFound, redirectLink)
	}

}

func (controller *SocialMediaLinkControllerImpl) GetLinkAnalytic(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)
	var request web.SocialMediaAnalyticGetRequest

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(err))
		return
	}

	socialMediaAnalyticResponse, errService := controller.AnalyticService.GetLinkAnalytic(ctx, request, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get social media link analytic",
			Data:    socialMediaAnalyticResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}

}

func (controller *SocialMediaLinkControllerImpl) GetSummaryLinkAnalytic(c *gin.Context) {
	ctx := context.Background()
	jwtToken := helper.ExtractTokenFromRequestHeader(c)

	socialMediaAnalyticSummaryResponse, errService := controller.AnalyticService.GetSummaryLinkAnalytic(ctx, jwtToken)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success get social media link analytic summary",
			Data:    socialMediaAnalyticSummaryResponse,
		}
		c.JSON(http.StatusOK, webResponse)
	}

}
