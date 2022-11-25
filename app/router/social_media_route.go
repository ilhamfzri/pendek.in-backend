package router

import (
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
	"gorm.io/gorm"
)

func AddSocialMediaRoute(server *Server, DB *gorm.DB, logger *logger.Logger, jwt helper.IJwt) {
	socialMediaTypeRepository := repository.NewSocialMediaTypeRepository(logger)
	socialMediaLinkRepository := repository.NewSocialMediaLinkRepository(logger)
	userRepository := repository.NewUserRepository(logger)

	socialMediaLinkService := service.NewSocialMediaLinkService(userRepository, socialMediaLinkRepository, socialMediaTypeRepository, DB, logger, jwt)
	socialMediaLinkController := controller.NewSocialMediaLink(socialMediaLinkService, logger)

	jwtMiddleware := middleware.NewJwtMiddleware(jwt.GetSigningKey())
	socialMediaRouteAuth := server.Router.Group("/v1/link/social-media")
	socialMediaRouteAuth.Use(jwtMiddleware)
	{
		socialMediaRouteAuth.GET("/types", socialMediaLinkController.GetAllTypes)
		socialMediaRouteAuth.POST("/", socialMediaLinkController.CreateLink)
		socialMediaRouteAuth.PUT("/:type_id", socialMediaLinkController.UpdateLink)
		socialMediaRouteAuth.GET("/", socialMediaLinkController.GetAllLink)
	}

	socialMediaRoute := server.Router.Group("")
	socialMediaRoute.GET("/:username/:social-media", socialMediaLinkController.RedirectLink)

}
