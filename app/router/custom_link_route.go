package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/handler"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
)

func AddCustomLinkRoute(server *Server, customLinkController controller.CustomLinkController, jwt helper.IJwt) {
	jwtMiddleware := middleware.NewJwtMiddleware(jwt.GetSigningKey())
	customLinkRouteAuth := server.Router.Group("/v1/link/custom")
	customLinkRouteAuth.Use(jwtMiddleware)
	{
		customLinkRouteAuth.POST("/", customLinkController.CreateLink)
		customLinkRouteAuth.GET("/", customLinkController.GetAllLink)
		customLinkRouteAuth.GET("/:link_id", customLinkController.GetLink)
		customLinkRouteAuth.PUT("/:link_id", customLinkController.UpdateLink)
		customLinkRouteAuth.POST("/upload-thumbnail", customLinkController.UploadCustomThumbnail)
		customLinkRouteAuth.GET("/user-thumbnail-list", customLinkController.GetUserThumbnail)
		customLinkRouteAuth.GET("/default-thumbnail-list", customLinkController.GetAllThumbnail)
		customLinkRouteAuth.GET("/check-short-code", customLinkController.CheckShortLinkAvaibility)
		customLinkRouteAuth.GET("/analytic", customLinkController.GetLinkAnalytic)
		customLinkRouteAuth.GET("/analytic/summary", customLinkController.GetSummaryLinkAnalytic)
	}

	customThumbnailResourcePath := os.Getenv("THUMBNAIL_IMG_DIR")
	customThumbnailResourceDirectory := http.Dir(customThumbnailResourcePath)
	customThumbnailRouteResources := server.Router.Group("/v1/resources/thumbnail")

	customThumbnailRouteResources.Use(func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasSuffix(urlPath, "/") {
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NoRouteResponse)
		}
		c.Next()
	})
	{
		customThumbnailRouteResources.StaticFS("/", customThumbnailResourceDirectory)
	}

	redirectCustomLinkRoute := server.Router.Group("/l")
	redirectCustomLinkRoute.GET("/:short_link_code", customLinkController.RedirectLink)
}
