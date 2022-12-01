package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/handler"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
	"gorm.io/gorm"
)

func AddCustomLinkRoute(server *Server, DB *gorm.DB, redis *redis.Client, logger *logger.Logger, jwt helper.IJwt) {
	customLinkRepository := repository.NewCustomLinkRepository(logger)
	customThumbnailRepository := repository.NewCustomThumbnailRepository(logger)
	thumbnailRepository := repository.NewThumbnailRepository(logger)

	customLinkService := service.NewCustomLinkService(customLinkRepository, customThumbnailRepository, thumbnailRepository, DB, logger, jwt)
	customLinkController := controller.NewCustomLinkController(customLinkService, logger)

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
